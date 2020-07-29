package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleSpriteCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	if len(env.args) == 0 {
		return botError{
			title:   "Validation Error",
			details: "Please enter a Pokémon name to get its sprite.",
		}
	}

	pkmArgs := parsePokemonCommand(env.args)

	// if the name and shininess were not parsed properly, lets assume it
	// follows the order on the help description.
	if pkmArgs.name == "" {
		pkmArgs.name = strings.ReplaceAll(env.args[0], "*", "")
		pkmArgs.isShiny = strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
	}

	pkm, err := b.pokemonRepo.pokemon(strings.ToLower(pkmArgs.name))
	if err != nil {
		return botError{
			title:   "Pokémon not found",
			details: fmt.Sprintf("Pokémon %s could not be found.", pkmArgs.name),
		}
	}

	var embedTitle string
	if pkmArgs.isShiny {
		embedTitle = "Shiny "
	}
	if pkmArgs.form != "" {
		embedTitle += strings.Title(pkmArgs.form) + " "
	}
	embedTitle += pkm.Name

	embed := b.newEmbed()
	embed.Title = embedTitle
	embed.Image = &discordgo.MessageEmbedImage{
		URL:    pkm.spriteImage(pkmArgs.isShiny, pkmArgs.form),
		Width:  300,
		Height: 300,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
