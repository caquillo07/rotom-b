package bot

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kubastick/dblgo"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/conf"
	"github.com/caquillo07/rotom-bot/repository"
)

type Bot struct {
	config         conf.Config
	session        *discordgo.Session
	repository     *repository.Repository
	commands       map[string]*command
	requestsServed uint64
}

func NewBot(conf conf.Config, session *discordgo.Session) *Bot {
	return &Bot{
		config:   conf,
		session:  session,
		commands: make(map[string]*command),
	}
}

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

	// Register ready as a callback for the ready events.
	b.session.AddHandler(b.ready)

	// Register the handleMessage func as a callback for MessageCreate events.
	b.session.AddHandler(b.handleMessage)

	// Open a websocket connection to Discord and begin listening.
	if err := b.session.Open(); err != nil {
		logger.Error("error opening connection", zap.Error(err))
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	logger.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	logger.Info("Shutting down...")
	return b.session.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	zap.L().Info("Bot is ready!")

	// Set the playing status.
	if err := s.UpdateStatus(0, b.config.Discord.PlayingStatus); err != nil {
		zap.L().Error("error setting bot playing status", zap.Error(err))
	}

	guilds, err := b.session.UserGuilds(0, "", "")
	if err != nil {
		zap.L().Error("error getting guild counts", zap.Error(err))
		return
	}
	zap.L().Info("Total guilds Rotom-B is on", zap.Int("guild_count", len(guilds)))

	if !b.config.DBL.Enable {
		return
	}

	dbl := dblgo.NewDBLApi(b.config.DBL.Token)
	err = dbl.PostStatsSimple(len(guilds))
	if err != nil {
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

	// Update the requests served, so we can get new ID for the next request
	reqID := atomic.AddUint64(&b.requestsServed, 1)
	logger.Info(
		"processing command",
		zap.Uint64("request_id", reqID),
		zap.String("command", m.Content),
		zap.String("user", m.Author.String()),
		zap.String("channel", channel.Name),
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
		botCmd = b.commands[botCmd.alias]
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
					m.Author.String() + "user doesn't have admin role",
					zap.String("command", cleanedMsg),
					zap.String("user", m.Author.String()),
				)
				return
			}
		}

		env := &commandEnvironment{
			args:    cmdParts[1:],
			command: cmdParts[0],
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
