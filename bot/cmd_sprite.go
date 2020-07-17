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
			details: "Please enter a Pokémon name to get its Pokédex info.",
		}
	}

	var embed *discordgo.MessageEmbed
	var err error

	isShiny := strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
	cleanPkmName := strings.ReplaceAll(env.args[0], "*", "")

	pkm, err := b.pokemonRepo.pokemon(strings.ToLower(cleanPkmName))
	if err != nil {
		return botError{
			title: "Pokémon not found",
			details: fmt.Sprintf("Pokémon %s could not be found.",
				cleanPkmName),
		}
	}

	pkmForm := ""
	if len(env.args) > 1 {
		pkmForm = getSpriteForm(env.args[1])
	}

	var embedTitle string
	if isShiny {
		embedTitle = "Shiny "
	}
	if pkmForm != "" {
		embedTitle += strings.Title(pkmForm) + " "
	}
	embedTitle += pkm.Name


	embed = b.newEmbed()
	embed.Title = embedTitle

	embed.Image = &discordgo.MessageEmbedImage{
		URL:    pkm.spriteImage(isShiny, pkmForm),
		Width:  300,
		Height: 300,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
