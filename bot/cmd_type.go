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

	offensiveMap := b.generateTypeMap(typeInfo.Offensive)
	defensiveMap := b.generateTypeMap(typeInfo.Defensive)

	offensiveNoDamageText := b.generateTypeText(offensiveMap[0])
	offensiveResistantText := b.generateTypeText(offensiveMap[0.5])
	offensiveSuperEffectiveText := b.generateTypeText(offensiveMap[2])

	defensiveNoDamageText := b.generateTypeText(defensiveMap[0])
	defensiveResistantText := b.generateTypeText(defensiveMap[0.5])
	defensiveSuperEffectiveText := b.generateTypeText(defensiveMap[2])

	embed = b.newEmbed()
	embed.Title = fmt.Sprintf("%s Type Info", strings.Title(pkmType))
	embed.Description = fmt.Sprintf(
		"%s Type Weakness, Resistances and Immunities", strings.Title(pkmType))
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

func (b *Bot) generateTypeMap(typeInfo map[string]float64) map[float64][]string {

	typesMap := make(map[float64][]string)

	for typeInInfo, num := range typeInfo {
		if num == 0 {
			typesMap[0] = append(typesMap[0], typeInInfo)
		}
		if num == 0.5 {
			typesMap[0.5] = append(typesMap[0.5], typeInInfo)
		}
		if num == 2 {
			typesMap[2] = append(typesMap[2], typeInInfo)
		}
	}

	return typesMap
}

func (b *Bot) generateTypeText(typesList []string) string {

	if len(typesList) == 0 {
		typesText := "None"
		return typesText
	}

	typesText := strings.Join(typesList, ", ") + "."

	return typesText
}
