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

	var embed *discordgo.MessageEmbed
	var err error

	pkmType := env.args[0]

	typeInfo, err := b.pokemonRepo.pokemonType(strings.ToLower(pkmType))
	if err != nil {
		return botError{
			title: "Type not found",
			details: fmt.Sprintf("Type %s could not be found.",
				pkmType),
		}
	}

	offensiveMap := make(map[float64][]string)
	defensiveMap := make(map[float64][]string)

	for offensiveType, num := range typeInfo.Offensive {
		if num == 0 {
			offensiveMap[0] = append(offensiveMap[0], offensiveType)
		}
		if num == 0.5 {
			offensiveMap[0.5] = append(offensiveMap[0.5], offensiveType)
		}
		if num == 2 {
			offensiveMap[2] = append(offensiveMap[2], offensiveType)
		}
	}

	for defensiveType, num := range typeInfo.Defensive {
		if num == 0 {
			defensiveMap[0] = append(defensiveMap[0], defensiveType)
		}
		if num == 0.5 {
			defensiveMap[0.5] = append(defensiveMap[0.5], defensiveType)
		}
		if num == 2 {
			defensiveMap[2] = append(defensiveMap[2], defensiveType)
		}
	}

	offensiveNoDamageText := b.generateTypeText(offensiveMap[0])
	offensiveResistantText := b.generateTypeText(offensiveMap[0.5])
	offensiveSuperEffectiveText := b.generateTypeText(offensiveMap[2])

	defensiveNoDamageText := b.generateTypeText(defensiveMap[0])
	defensiveResistantText := b.generateTypeText(defensiveMap[0.5])
	defensiveSuperEffectiveText := b.generateTypeText(defensiveMap[2])

	embed = b.newEmbed()
	embed.Title = fmt.Sprintf("%s Type Info", strings.Title(pkmType))
	embed.Description = fmt.Sprintf("%s Type Weakness, Resistances and Immunities", strings.Title(pkmType))
	embed.Color = typeInfo.Color
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name: "Offensive",
			Value: fmt.Sprintf("Super Effective: `%s`\nResistant: `%s`\nNo Damage: `%s`",
				offensiveSuperEffectiveText,
				offensiveResistantText,
				offensiveNoDamageText,
			),
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name: "Defensive",
			Value: fmt.Sprintf(
				"Super Effective: `%s`\nResistant: `%s`\nNo Damage: `%s`",
				defensiveSuperEffectiveText,
				defensiveResistantText,
				defensiveNoDamageText,
			),
			Inline: false,
		},
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err

}

func (b *Bot) generateTypeText(typesList []string) string {

	if len(typesList) == 0 {
		typesText := "None"
		return typesText
	}

	var typesText string

	for i, typeInList := range typesList {
		if i == len(typesList)-1 {
			typesText += typeInList + "."
		} else {
			typesText += typeInList + ", "
		}
	}

	return typesText
}
