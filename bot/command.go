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
		execute: b.handleCatchCmd,
		helpText: `Shows a detailed summary of catch rates for a given Pokémon and Ball combination.

This command will perform the calculations as presented by [Bulbapedia](https://bulbapedia.bulbagarden.net/wiki/Catch_rate#Probability_of_capture). 
The calculations are estimates, and due to a [rounding error](https://bulbapedia.bulbagarden.net/wiki/Catch_rate#Probability_of_capture), at some points its impossible to calculate with accuracy.

The confidence level will display when this calculations fall under the rounding error`,
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
		helpText:  "Shows ithe Pokémon Sprite in appropriate form",
		usage:     b.addCmdPrefix("{{p}}nature <nature>"),
		example:   b.addCmdPrefix("{{p}}nature modest"),
		adminOnly: false,
	}
	b.commands["type"] = &command{
		execute:   b.handleTypeCmd,
		helpText:  "Shows info regarding Pokémon Types.",
		usage:     b.addCmdPrefix("{{p}}type <type>"),
		example:   b.addCmdPrefix("{{p}}type grass"),
		adminOnly: false,
	}
	b.commands["pokedex"] = &command{
		execute:   b.handlePokedexCmd,
		helpText:  "Shows Pokédex info on every Pokémon.",
		usage:     b.addCmdPrefix("{{p}}pokedex <pokemon>"),
		example:   b.addCmdPrefix("{{p}}pokedex caterpie"),
		adminOnly: false,
	}
	b.commands["sprite"] = &command{
		execute:   b.handleSpriteCmd,
		helpText:  "Shows the Pokémon Sprite. Include * in the end for the shiny sprite.",
		usage:     b.addCmdPrefix("{{p}}sprite <pokemon>"),
		example:   b.addCmdPrefix("{{p}}sprite charizard* gmax"),
		adminOnly: false,
	}
	b.commands["invite"] = &command{
		execute:   b.handleInviteCmd,
		helpText:  "Get an Invite Link to invite Rotom-B to another server!",
		usage:     b.addCmdPrefix("{{p}}invite"),
		example:   b.addCmdPrefix("{{p}}invite"),
		adminOnly: false,
	}
	b.commands["version"] = &command{
		execute:   b.handleVersionCmd,
		helpText:  "Check which version of Rotom-B is running.",
		usage:     b.addCmdPrefix("{{p}}version"),
		example:   b.addCmdPrefix("{{p}}version"),
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
