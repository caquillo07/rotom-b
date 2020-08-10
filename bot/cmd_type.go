package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleTypeCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	if len(env.args) == 0 {
		return botError{
			title:   "Validation Error",
			details: "Please enter a den number or a Pokémon name to look for related dens.",
		}
	}

	pkmType := env.args[0]
	typeInfo, err := b.repository.PokemonType(strings.ToLower(pkmType))
	if err != nil {
		return botError{
			title: "Type not found",
			details: fmt.Sprintf("Type %s could not be found.",
				pkmType),
		}
	}

	embed := b.newEmbed()
	embed.Title = fmt.Sprintf("%s Type Info", strings.Title(pkmType))
	embed.Description = fmt.Sprintf(
		"%s Type Weakness, Resistances and Immunities",
		strings.Title(pkmType),
	)
	embed.Color = typeInfo.Color
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Offensive",
			Value:  generateTypeMessage(typeInfo.Offensive),
			Inline: false,
		},
		{
			Name:   "Defensive",
			Value:  generateTypeMessage(typeInfo.Defensive),
			Inline: false,
		},
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func generateTypeMessage(typeInfo map[string]float64) string {
	typesMap := make(map[float64][]string)
	for t, num := range typeInfo {
		typesMap[num] = append(typesMap[num], t)
	}

	return fmt.Sprintf(
		"Super Effective: `%s`\nResistant: `%s`\nNo Damage: `%s`",
		generateTypeText(typesMap[2]),
		generateTypeText(typesMap[0.5]),
		generateTypeText(typesMap[0]),
	)
}

func generateTypeText(typesList []string) string {
	if len(typesList) == 0 {
		return "None"
	}
	return strings.Join(typesList, ", ") + "."
}
