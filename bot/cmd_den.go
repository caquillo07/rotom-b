package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleDenCmd(
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
	_, err := strconv.Atoi(env.args[0])
	isNumber := err == nil

	if isNumber {
		embed, err = b.getDenFromNumber(env.args[0])
	} else {
		isShiny := strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
		cleanPkmName := strings.ReplaceAll(env.args[0], "*", "")
		form := getFormFromArgs(env.args)
		embed, err = b.getDensFromPokemon(cleanPkmName, form, isShiny)
	}
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func (b *Bot) getDensFromPokemon(pkmnName, form string, isShiny bool) (*discordgo.MessageEmbed, error) {

	pkm, err := b.pokemonRepo.pokemon(strings.ToLower(pkmnName))
	if err != nil {
		return nil, botError{
			title: "Pokémon not found",
			details: fmt.Sprintf("Pokémon %s could not be found.",
				pkmnName),
		}
	}

	embed := b.newEmbed()
	embed.Title = pkm.Name + " is in the following Dens:"
	embed.Image = &discordgo.MessageEmbedImage{
		URL:    pkm.spriteImage(isShiny, form),
		Width:  300,
		Height: 300,
	}

	if txt := getDensText(pkm.Dens.Sword); txt != "" {
		if txt := getDensText(pkm.Dens.Shield); txt != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Inline: true,
				Name:   "*Sword:* ",
				Value:  txt,
			})
		}
	}

	if txt := getDensText(pkm.Dens.Shield); txt != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Inline: true,
			Name:   "*Shield:* ",
			Value:  txt,
		})
	}

	return embed, nil
}

func getDensText(dens []string) string {
	var txt string
	for i := 0; i < len(dens); i++ {
		den := strings.ToLower(dens[i])
		txt += fmt.Sprintf(
			"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml)",
			den,
			den,
		)
		if i != len(dens)-1 {
			txt += ", "
		}
	}
	return txt
}

func (b *Bot) getDenFromNumber(denNumber string) (*discordgo.MessageEmbed, error) {

	den, err := b.pokemonRepo.den(denNumber)
	if err != nil {
		return nil, botError{
			title: "Den number not found",
			details: fmt.Sprintf("Den %s could not be found.",
				denNumber),
		}
	}

	swordField := &discordgo.MessageEmbedField{}
	swordField.Inline = true
	swordField.Name += "HA in Sword"
	for i := 0; i < len(den.Sword); i++ {
		if den.Sword[i].Ability != "Standard" {
			swordField.Value += den.Sword[i].Name + "\n"
		}
	}

	shieldField := &discordgo.MessageEmbedField{}
	shieldField.Inline = true
	shieldField.Name += "HA in Shield"
	for i := 0; i < len(den.Shield); i++ {
		if den.Shield[i].Ability != "Standard" {
			shieldField.Value += den.Shield[i].Name + "\n"
		}
	}

	embed := b.newEmbed()
	embed.Title = "Pokémon found in Den " + den.Number + ":"
	embed.URL = fmt.Sprintf(
		"https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml",
		strings.ToLower(strings.ReplaceAll(den.Number, " ", "")),
	)
	embed.Image = &discordgo.MessageEmbedImage{
		URL: fmt.Sprintf(
			"https://caquillo07.github.io/data/dens/den_%s.png",
			strings.ToLower(strings.ReplaceAll(den.Number, " ", "")),
		),
	}
	embed.Fields = []*discordgo.MessageEmbedField{swordField, shieldField}
	return embed, nil
}
