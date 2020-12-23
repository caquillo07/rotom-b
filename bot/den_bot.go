package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kubastick/dblgo"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/conf"
	"github.com/caquillo07/rotom-bot/repository"
)

const galarian = "galarian"

// Bot is the struct defining the den bot, it is responsible for listening
// to on the discord session and handling messages.
type Bot struct {
	config         conf.Config
	sessions       []*discordgo.Session
	repository     *repository.Repository
	commands       map[string]*command
	requestsServed uint64
}

// NewBot creates a new bot instance from the given session and config
func NewBot(conf conf.Config) *Bot {
	return &Bot{
		config:   conf,
		commands: make(map[string]*command),
	}
}

// Run starts an instance of a bot. This method will create a new repository
// with a database connection to the given configuration, initialize all
// commands and listen on the discord web socket.
func (b *Bot) Run() error {
	logger := zap.L()

	// Open connection to DB before anything else
	db, err := repository.Open(b.config.Database)
	if err != nil {
		return err
	}

	// Initialize the repo
	repo, err := repository.NewRepository(db)
	if err != nil {
		return err
	}
	b.repository = repo

	// Create all the commands on the bot
	b.initCommands()

	// Create a new Discord session using the provided login information.
	// This session will help us get the recommended sharding numbers.
	discordToken := "Bot " + b.config.Discord.Token
	mainSession, err := discordgo.New(discordToken)
	if err != nil {
		logger.Fatal("error creating Discord session", zap.Error(err))
	}

	gateway, err := mainSession.GatewayBot()
	if err != nil {
		return errors.Wrap(err, "failed to create gateway")
	}

	logger.Info(
		fmt.Sprintf("opening bot with %d shards", gateway.Shards),
		zap.Int("shards", gateway.Shards),
	)

	// Now lets create a new session for each of the shards we get.
	b.sessions = make([]*discordgo.Session, gateway.Shards)
	wg := sync.WaitGroup{}
	for i := 0; i < gateway.Shards; i++ {
		logger.Info(fmt.Sprintf("opening shared %d", i+1))
		session, err := discordgo.New(discordToken)
		if err != nil {
			return err
		}

		session.ShardCount = gateway.Shards
		session.ShardID = i
		session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
		session.State.TrackPresences = false
		b.sessions[i] = session
		wg.Add(1)

		// We shoot the connection off on its own to speed it up a bit.
		go func(session *discordgo.Session, shard int) {
			// Open a websocket connection to Discord and begin listening.
			defer wg.Done()
			if err := session.Open(); err != nil {
				logger.Error(fmt.Sprintf("error opening connection on shard %d", shard), zap.Error(err))
			}
		}(b.sessions[i], i+1)
	}
	wg.Wait()

	// Add all handlers
	for _, session := range b.sessions {

		// Register ready as a callback for the ready events.
		session.AddHandler(b.ready)

		// Register the handleMessage func as a callback for MessageCreate events.
		session.AddHandler(b.handleMessage)
	}

	// Wait here until CTRL-C or other term signal is received.
	logger.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	logger.Info("Shutting down...")
	var closingError error
	for i, session := range b.sessions {

		if err := session.Close(); err != nil {
			closingError = errors.Wrap(err, fmt.Sprintf("failed to close shard %d", i))
		}
	}
	return closingError
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	zap.L().Info("Bot is ready!")

	// Set the playing status.
	if err := s.UpdateStatus(0, b.config.Discord.PlayingStatus); err != nil {
		zap.L().Error("error setting bot playing status", zap.Error(err))
	}

	guilds := len(s.State.Guilds)
	zap.L().Info("Total guilds Rotom-B is on", zap.Int("guild_count", guilds))

	if !b.config.DBL.Enable {
		return
	}

	if err := dblgo.NewDBLApi(b.config.DBL.Token).PostStatsSimple(guilds); err != nil {
		zap.L().Error("error posting user count to DBL", zap.Error(err))
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself, or any other bot for that
	// matter
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	// Lets start by fetching the guild's config and and checking if the message
	// is for us.
	logger := zap.L()
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		logger.Error(
			"error retrieving the channel",
			zap.Error(err),
			zap.String("channel_id", channel.ID),
		)
		return
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		logger.Error(
			"error retrieving the guild",
			zap.Error(err),
			zap.String("guild_id", channel.GuildID),
		)
		return
	}

	guildSettings, err := b.getOrCreateGuildSettings(guild)
	if err != nil {
		logger.Error(
			"error retrieving the guild's config",
			zap.Error(err),
			zap.String("guild_name", guild.Name),
			zap.String("guild_id", guild.ID),
		)
	}

	prefix := b.config.Bot.Prefix
	if pre := guildSettings.BotPrefix; pre != "" {
		prefix = pre
	}

	// If the message does not have the required prefix, exit as well
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// If this guild has specific channels to listen on and this channel is not
	// in it, exit.
	if len(guildSettings.ListeningChannels) > 0 && !inListenChannels(channel.ID, guildSettings.ListeningChannels) {
		return
	}

	// Update the requests served, so we can get new ID for the next request
	reqID := atomic.AddUint64(&b.requestsServed, 1)
	logger.Info(
		"processing command",
		zap.Uint64("request_id", reqID),
		zap.String("command", m.Content),
		zap.String("user", m.Author.String()),
		zap.String("channel", channel.Name),
		zap.String("guild_name", guild.Name),
	)

	// Handle panics gracefully, sucks that we do it this late but we need some
	// of the info gathered above. Special care needed when changing code above.
	defer func(reqID uint64) {
		if r := recover(); r != nil {
			b.handlePanic(r, s, m, reqID)
		}
	}(reqID)

	// To not have to check for the prefix on every single command
	cleanedMsg := strings.TrimPrefix(m.Content, prefix)
	cmdParts := strings.Split(cleanedMsg, " ")
	botCmd, ok := b.commands[cmdParts[0]]
	if !ok {
		// ignoring unknown commands
		return
	}
	if botCmd.alias != "" {
		botCmd, ok = b.commands[botCmd.alias]
		if !ok {
			return
		}
	}

	// Send this off on its own go routine to be able to handle many of them
	// at once
	go func(reqID uint64) {
		// before anything else, lets make sure the user has permissions
		// to perform this command
		if botCmd.adminOnly {
			isAdmin, err := userIsAdmin(s, guild.ID, m.Author.ID)
			if err != nil {
				logger.Error(
					"failed to check user's admin status",
					zap.String("command", cleanedMsg),
					zap.Uint64("request_id", reqID),
					zap.Error(err),
				)

				b.handleCommandError(s, m, reqID, err)
				return
			}
			if !isAdmin {
				logger.Info(
					m.Author.String()+" user doesn't have admin role",
					zap.String("command", cleanedMsg),
					zap.String("user", m.Author.String()),
				)
				return
			}
		}

		env := &commandEnvironment{
			args:          cmdParts[1:],
			command:       cmdParts[0],
			commandPrefix: prefix,
		}
		if err := botCmd.execute(s, env, m.Message); err != nil {
			logger.Error(
				"failed to handle command",
				zap.String("command", cleanedMsg),
				zap.Uint64("request_id", reqID),
				zap.Error(err),
			)

			b.handleCommandError(s, m, reqID, err)
		}
	}(reqID)
}

func (b *Bot) handleCommandError(s *discordgo.Session, m *discordgo.MessageCreate, reqID uint64, err error) {
	// If the error is a botError, that means that we will consider
	// it public and just pass the error on to the user.
	errTitle := "Internal Error"
	errDetails := fmt.Sprintf(`Whoops, there was an error processing the request with ID **%d**`, reqID)
	if publicErr, ok := err.(botError); ok {
		errTitle, errDetails = publicErr.title, publicErr.details
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, b.newErrorEmbedf(errTitle, errDetails))

	// If this errors, then ¯\_(ツ)_/¯ log and move on
	if err != nil {
		zap.L().Error(
			"failed to communicate command error",
			zap.Error(err),
		)
	}
}

func (b *Bot) handlePanic(panic interface{}, s *discordgo.Session, m *discordgo.MessageCreate, reqID uint64) {
	logger := zap.L()
	logger.Error(
		"recovered from panic while handling message",
		zap.String("panic_message", fmt.Sprintf("%s", panic)),
		zap.String("command", m.Content),
		zap.String("user", m.Author.Username),
	)

	// Log the stacktrace to the console
	logger.Info(
		"recovered from panic",
		zap.String("panic_message", fmt.Sprintf("%s", panic)),
		zap.Stack("stack_trace"),
	)

	b.handleCommandError(s, m, reqID, errors.New("internal error"))
}

func (b *Bot) getOrCreateGuildSettings(guild *discordgo.Guild) (*repository.GuildSettings, error) {
	gc, err := b.repository.GetGuildSettings(guild.ID)
	if err != nil && err != repository.ErrRecordNotFound {
		return nil, err
	}

	if gc != nil {
		// TODO: todo so its highlighted, remove this once this is no longer a
		//  thing. If the name is <replace>, replace it with the real name.
		//  Since this is a temp action, ignore the error, just log it
		if gc.Name == "<replace>" {
			gc.Name = guild.Name
			if err := b.repository.UpdateGuildSettings(gc); err != nil {
				zap.L().Error(
					"failed to update guild name",
					zap.Error(err),
					zap.String("guild_id", guild.ID),
					zap.String("guild_name", guild.Name),
				)
			}
		}
		return gc, nil
	}

	// not found, then create
	gc = &repository.GuildSettings{
		Name:      guild.Name,
		DiscordID: guild.ID,
		BotPrefix: b.config.Bot.Prefix,
	}
	if err := b.repository.CreateGuildSettings(gc); err != nil {
		return nil, err
	}
	return gc, nil
}

func userIsAdmin(
	s *discordgo.Session,
	guildID, userID string,
) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		// ok to ignore this error since its just a not found error on the
		// local state
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true, nil
		}
	}
	return false, nil
}

func sendEmbed(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) error {
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func inListenChannels(id string, s []*repository.GuildSettingChannel) bool {
	for _, ss := range s {
		if id == ss.ID {
			return true
		}
	}
	return false
}
