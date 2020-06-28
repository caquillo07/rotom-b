package bot

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/conf"
)

type Bot struct {
	config      conf.Config
	session     *discordgo.Session
	pokemonRepo *pokemonRepo
}

func NewBot(conf conf.Config, session *discordgo.Session) *Bot {
	return &Bot{
		config:  conf,
		session: session,
	}
}

func (b *Bot) Run() error {
	logger := zap.L()

	// Initialize the repo
	repo, err := newPokemonRepo()
	if err != nil {
		return err
	}
	b.pokemonRepo = repo

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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message does not have the required prefix, exit as well
	if !strings.HasPrefix(m.Content, b.config.Bot.Prefix) {
		return
	}

	// To not have to check for the prefix on every single command
	command := strings.TrimPrefix(m.Content, b.config.Bot.Prefix)
	logger := zap.L()
	switch command {
	case "ping":
		if err := b.handlePingCmd(s, m); err != nil {
			logger.Error("failed to respond to ping command", zap.Error(err))
		}
	case "pong":
		if err := handlePongCmd(s, m); err != nil {
			logger.Error("failed to respond to ping command", zap.Error(err))
		}
	case "den":
		if err := b.handleDenCmd(s, m); err != nil {
			logger.Error("failed to respond to den command", zap.Error(err))
		}
	default:
		// ignoring unknown commands
	}
}
