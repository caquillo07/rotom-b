package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type command struct {
	execute   func(s *discordgo.Session, env *commandEnvironment, m *discordgo.Message) error
	helpText  string
	usage     string
	example   string
	adminOnly bool
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
	b.commands["ball"] = &command{
		execute:   b.handleBallCmd,
		helpText:  "Shows a summary of a Poké-Ball’s statistics",
		usage:     b.addCmdPrefix("{{p}}ball <ball_name>"),
		example:   b.addCmdPrefix("{{p}}ball beast"),
		adminOnly: false,
	}
	b.commands["catch"] = &command{
		execute:   b.handleCatchCmd,
		helpText:  "Shows a detailed summary of catch rates for a given Pokémon and Ball combination.",
		usage:     b.addCmdPrefix("{{p}}catch <pokemon> [form] [ball_name]"),
		example:   b.addCmdPrefix("{{p}}catch charizard gmax lux"),
		adminOnly: false,
	}
	b.commands["credits"] = &command{
		execute:   b.handleCreditsCmd,
		helpText:  "Credits to all who helped in the creation of the bot.",
		usage:     b.addCmdPrefix("{{p}}credits"),
		example:   b.addCmdPrefix("{{p}}credits"),
		adminOnly: false,
	}
	b.commands["nature"] = &command{
		execute:   b.handleNatureCmd,
		helpText:  "Shows info regarding Pokémon Natures.",
		usage:     b.addCmdPrefix("{{p}}nature <nature>"),
		example:   b.addCmdPrefix("{{p}}nature modest"),
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

func (b *Bot) newErrorEmbedf(errorTitle, errorMsg string, a ...interface{}) *discordgo.MessageEmbed {
	embed := b.newEmbed()
	embed.Title = errorTitle
	embed.Description = fmt.Sprintf(errorMsg, a...)
	embed.Color = b.config.Bot.ErrorEmbedColor
	return embed
}
