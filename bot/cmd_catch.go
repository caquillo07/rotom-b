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
	form := getFormFromArgs(env.args)
	ball, err := b.pokemonRepo.ball(env.args[len(env.args)-1])
	if err != errBallDoesNotExist {
		return err
	}

	// now make sure that someone did not just send in a ball and no pokemon
	if ball != nil && len(env.args) == 1 {
		return botError{
			title:   "Validation Error",
			details: "A pokemon must be provided alongside the PokeBall for the catch command",
		}
	}

	// If the ball does not exist, that means we got just a pokemon request
	if ball == nil {
		embed, err := b.getPokemonTopFourBalls(cleanPkmName, form, isShiny)
		if err != nil {
			return err
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	}

	// If we got a ball, we doing an specific check against a pokemon. First
	// sure we did get a pokemon
	_, err = s.ChannelMessage(m.ChannelID, cleanPkmName+" "+ball.Name)
	return err
}

func (b *Bot) getPokemonTopFourBalls(pokemonName, form string, shiny bool) (*discordgo.MessageEmbed, error) {
	pkm, err := b.pokemonRepo.pokemon(pokemonName)
	if err != nil {
		return nil, botError{
			title:   "Not Found",
			details: fmt.Sprintf("Pokemon %s was not found", pokemonName),
		}
	}

	form = getSpriteForm(form)
	name := pkm.Name
	if form == gigantamax {
		name = "G-Max " + name
	}

	embed := b.newEmbed()
	embed.Title = "Best Catch Rates"
	embed.Description = fmt.Sprintf("The best balls for catching %s are:", name)
	// todo(hector) add support for color

	embed.Image = &discordgo.MessageEmbedImage{
		URL: pkm.spriteImage(shiny, form),
	}

	return embed, nil
}

func getFormFromArgs(args []string) string {

	// the form will always be the second element
	// e.g: blastoise gmax
	if len(args) < 2 || args[1] == "" {
		return ""
	}

	return getSpriteForm(args[1])
}
