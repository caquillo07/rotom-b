package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
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
	pkmType := env.args[0]
	typeInfo, err := b.pokemonRepo.pokemonType(strings.ToLower(pkmType))
	if err != nil {
		return botError{
			title: "Type not found",
			details: fmt.Sprintf("Type %s could not be found.",
				pkmType),
		}
	}

	var (
		offensiveSuperEffective []string
		offensiveResistant      []string
		offensiveNoDamage       []string
		defensiveSuperEffective []string
		defensiveResistant      []string
		defensiveNoDamage       []string
	)

	for offensiveType, num := range typeInfo.Offensive {
		if num == 0 {
			offensiveNoDamage = append(offensiveNoDamage, offensiveType)
		}
		if num == 0.5 {
			offensiveResistant = append(offensiveResistant, offensiveType)
		}
		if num == 2 {
			offensiveSuperEffective = append(offensiveSuperEffective, offensiveType)
		}
	}

	for defensiveType, num := range typeInfo.Defensive {
		if num == 0 {
			defensiveNoDamage = append(defensiveNoDamage, defensiveType)
		}
		if num == 0.5 {
			defensiveResistant = append(defensiveResistant, defensiveType)
		}
		if num == 2 {
			defensiveSuperEffective = append(defensiveSuperEffective, defensiveType)
		}
	}

	offensiveNoDamageText := "None"
	for i, offensiveType := range offensiveNoDamage {
		if i == len(offensiveNoDamage)-1 {
			offensiveNoDamageText += offensiveType + "."
		}
		offensiveNoDamageText += offensiveType + ", "
	}

	var offensiveResistantText string
	for i, offensiveType := range offensiveResistant {
		if i == len(offensiveResistant)-1 {
			offensiveResistantText += offensiveType + "."
		}
		offensiveResistantText += offensiveType + ", "
	}

	var offensiveSuperEffectiveText string
	for i, offensiveType := range offensiveSuperEffective {
		if i == len(offensiveSuperEffective)-1 {
			offensiveSuperEffectiveText += offensiveType + "."
		}
		offensiveSuperEffectiveText += offensiveType + ", "
	}

	defensiveNoDamageText := "None"
	for i, defensiveType := range defensiveNoDamage {
		if i == len(defensiveNoDamage)-1 {
			defensiveNoDamageText += defensiveType + "."
		}
		defensiveNoDamageText += defensiveType + ", "
	}

	var defensiveResistantText string
	for i, defensiveType := range defensiveResistant {
		if i == len(defensiveResistant)-1 {
			defensiveResistantText += defensiveType + "."
		}
		defensiveResistantText += defensiveType + ", "
	}

	var defensiveSuperEffectiveText string
	for i, defensiveType := range defensiveSuperEffective {
		if i == len(defensiveSuperEffective)-1 {
			defensiveSuperEffectiveText += defensiveType + "."
		}
		defensiveSuperEffectiveText += defensiveType + ", "
	}

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
