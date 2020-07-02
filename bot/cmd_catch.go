package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// handleCatchCmd handles the catch command, sends back a
// detailed summary of catch rates for a given Pokémon & Ball.
func (b *Bot) handleCatchCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {
	if len(env.args) == 0 {
		return botError{
			title:   "Validation Error",
			details: "Please enter a Pokémon to catch followed by a Poké-Ball of your choice.",
		}
	}

	isShiny := strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
	cleanPkmName := strings.ReplaceAll(env.args[0], "*", "")
	var embed *discordgo.MessageEmbed
	var err error

	// we just got a pokemon request
	if len(env.args) == 1 {
		embed, err = b.getPokemonTopFourBalls(cleanPkmName, isShiny)
		if err != nil {
			return err
		}
	} else if len(env.args) == 2 {
		// we got a pokemon and ball request

	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func (b *Bot) getPokemonTopFourBalls(pokemonName string, shiny bool) (*discordgo.MessageEmbed, error) {
	pkm, err := b.pokemonRepo.pokemon(pokemonName)
	if err != nil {
		return nil, botError{
			title:   "Not Found",
			details: fmt.Sprintf("%s was not found", pokemonName),
		}
	}

	name := pkm.Name
	if pkm.isGigantamax() {
		name = "G-Max " + name
	}

	embed := b.newEmbed()
	embed.Title = "Best Catch Rates"
	embed.Description = fmt.Sprintf("The best balls for catching G-Max %s are:", name)
	// todo(hector) add support for color

	embed.Image = &discordgo.MessageEmbedImage{
		// todo(hector) add form
		URL: pkm.spriteImage(shiny, ""),
	}

	return embed, nil
}
