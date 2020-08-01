package bot

import (
	"fmt"
	"strconv"
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

type pokemonArg struct {
	name    string
	isShiny bool
	ball    string
	form    string
	den     string
	extras  []string
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

func parsePokemonCommand(args []string) pokemonArg {

	// make a uniformed string to make it easier to parse
	for i, a := range args {
		args[i] = strings.ToLower(a)
	}

	// then test each thing we are looking for, this is not optimal but it works
	// for now.
	pkmArgs := pokemonArg{}
	var skipIndex bool
	for i, arg := range args {

		// in some instances, we already processed the this index on the
		// previous index. In such cases, simply move on. Make sure to reset
		// the flag so we don't skip parsing the rest of the arguments.
		if skipIndex {
			skipIndex = false
			continue
		}

		// first we will try to find the form by looping through the arguments
		// and see if any of them is a valid form.
		if f := getSpriteForm(arg); f != "" {
			pkmArgs.form = f
			continue
		}

		// next, try to find a den number if any
		if _, err := strconv.Atoi(arg); err == nil {
			pkmArgs.den = arg
			continue
		}

		// now check if its a pokemon with a two-part name and its not the last
		// element in the slice
		cleanArg := strings.ReplaceAll(arg, "*", "")
		if (cleanArg == "mr" || cleanArg == "mr." || cleanArg == "mime") && i != len(args)-1 {
			secondPart := args[i+1]
			cleanSecondPart := strings.ReplaceAll(secondPart, "*", "")
			if strings.Contains(arg, "mr") && (cleanSecondPart == "rime" || cleanSecondPart == "mime") {
				pkmArgs.name = "mr " + cleanSecondPart
				pkmArgs.isShiny = strings.Contains(arg, "*") || strings.Contains(secondPart, "*")
				skipIndex = true
				continue
			}

			if cleanArg == "mime" && strings.Contains(cleanSecondPart, "jr") {
				pkmArgs.name = "mime jr"
				pkmArgs.isShiny = strings.Contains(arg, "*") || strings.Contains(secondPart, "*")
				skipIndex = true
				continue
			}

			// if for some reason we made it this far, let rest of the parser
			// continue. The leftover will just be added to the extras at the
			// end
		}

		// we have to check for special hyphenated names because the sprites
		// are named stupidly and inconsistently
		var found bool
		for _, s := range []string{"mr-mime", "mr.-mime", "mr-rime", "mr.-rime", "mime-jr", "mime-jr."} {
			if cleanArg != s {
				continue // continue this inner loop
			}

			// found a special one, process as needed, then break
			found = true
			s = strings.ReplaceAll(s, "-", " ")
			pkmArgs.name = strings.ReplaceAll(s, ".", "")
			pkmArgs.isShiny = strings.Contains(arg, "*")
			break
		}
		if found {
			continue
		}

		// we also have to check for special ' names because the sprites
		// are named stupidly and inconsistently
		for _, s := range []string{"farfetchd", "sirfetchd", "hooh"} {
			if cleanArg != s {
				continue // continue this inner loop
			}

			// found a special one, process as needed, then break
			found = true
			if s == "hooh" {
				pkmArgs.name = "ho-oh"
				pkmArgs.isShiny = strings.Contains(arg, "*")
				break
			}

			// in this case, both pokemon end with a 'd' and its the only 'd' in
			// the name. Lets just replace the 'd' and move on. This will NOT
			// work if a new pokemon with a ' is added
			pkmArgs.name = strings.ReplaceAll(s, "d", "'d")
			pkmArgs.isShiny = strings.Contains(arg, "*")
			break
		}
		if found {
			continue
		}

		// last, lets see if we can find a den on the command. This is a "heavy"
		// check since we are loop through an array inside another loop, so lets
		// do it last.
		for _, b := range ballNames {
			if strings.HasPrefix(arg, b) {
				pkmArgs.ball = arg
				found = true
				break // we done, break the inner loop
			}
		}
		if found {
			continue
		}

		// if we made it this far, none of the parsers caught it. Normally the
		// first argument is assumed to be the Pokemon's name, so just add it
		// to the name, and everything else add it to the extra/unknown field
		if pkmArgs.name == "" {
			pkmArgs.name = cleanArg
			pkmArgs.isShiny = strings.Contains(arg, "*")
			continue
		}
		pkmArgs.extras = append(pkmArgs.extras, arg)
	}

	return pkmArgs
}
