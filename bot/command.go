package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type command struct {
	execute   func(s *discordgo.Session, env *commandEnvironment, m *discordgo.Message) error
	helpText  string
	usage     string
	example   string
	adminOnly bool
	// requiredPermission int
}

type commandEnvironment struct {
	command string
	args    []string
}

func (b *Bot) initCommands() {
	b.commands["help"] = &command{
		execute:   b.handleHelpCmd,
		helpText:  "Displays a list of commands you have access to use.",
		usage:     b.addCmdPrefix("{{p}}help [command]"),
		example:   b.addCmdPrefix("{{p}}help\n{{p}}help den"),
		adminOnly: false,
	}
	b.commands["den"] = &command{
		execute:   b.handleDenCmd,
		helpText:  "Shows a list of Pokémon that belong to a den including their HAs.",
		usage:     b.addCmdPrefix("{{p}}den <den_number|pokemon_name>"),
		example:   b.addCmdPrefix("{{p}}den 22\n{{p}}den charizard"),
		adminOnly: false,
	}
}

// addCmdPrefix replaces all cases of {{p}} with the actual
// bot command prefix.
func (b *Bot) addCmdPrefix(s string) string {
	return strings.ReplaceAll(s, "{{p}}", b.config.Bot.Prefix)
}

func (b *Bot) newEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: b.config.Bot.EmbedColor,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Rotom-B - By Hector & Milla",
			IconURL: "https://images-na.ssl-images-amazon.com/images/I/41x0Y9yJYKL.jpg",
		},
	}
}
