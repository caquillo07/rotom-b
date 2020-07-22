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

	if len(pkm.Dens.Shield) == 0 && len(pkm.Dens.Sword) == 0 {
		return nil, botError{
			title: fmt.Sprintf(
				"%s is not in any Dens!",
				pkm.Name,
			),
			details: fmt.Sprintf(
				"%s could not be found in any current dens.",
				pkm.Name,
			),
		}
	}

	// TODO: Optimize this later, too much repeated code!

	//Sword
	var txtSword, txtHASword string
	for i := 0; i < len(pkm.Dens.Sword); i++ {
		den, _ := b.pokemonRepo.den(pkm.Dens.Sword[i])
		for j := 0; j < len(den.Sword); j++ {
			if strings.ToLower(den.Sword[j].Name) == strings.ToLower(pkm.Name) {
				if den.Sword[j].Ability == "Standard" {
					txtSword += fmt.Sprintf(
						"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml), ",
						den.Number,
						den.Number,
					)
				} else {
					txtHASword += fmt.Sprintf(
						"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml), ",
						den.Number,
						den.Number,
					)
				}
			}
		}
	}

	//Shield
	var txtShield, txtHAShield string
	for i := 0; i < len(pkm.Dens.Shield); i++ {
		den, _ := b.pokemonRepo.den(pkm.Dens.Shield[i])
		for j := 0; j < len(den.Shield); j++ {
			if strings.ToLower(den.Shield[j].Name) == strings.ToLower(pkm.Name) {
				if den.Shield[j].Ability == "Standard" {
					txtShield += fmt.Sprintf(
						"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml), ",
						den.Number,
						den.Number,
					)
				} else {
					txtHAShield += fmt.Sprintf(
						"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml), ",
						den.Number,
						den.Number,
					)
				}
			}
		}
	}

	embed := b.newEmbed()
	embed.Title = pkm.Name + " is in the following Dens:"
	embed.Image = &discordgo.MessageEmbedImage{
		URL:    pkm.spriteImage(isShiny, form),
		Width:  300,
		Height: 300,
	}
	embed.Fields = []*discordgo.MessageEmbedField{}

	if txtSword != "" || txtHASword != "" {
		if txtHASword != "" {
			txtHASword = "HA: " + txtHASword
		}
		densSword := &discordgo.MessageEmbedField{
			Name:  "Sword",
			Value: txtSword + txtHASword,
		}
		embed.Fields = append(embed.Fields, densSword)
	}

	if txtShield != "" || txtHAShield != "" {
		if txtHAShield != "" {
			txtHAShield = "HA: " + txtHAShield
		}
		densShield := &discordgo.MessageEmbedField{
			Name:  "Sword",
			Value: txtShield + txtHAShield,
		}
		embed.Fields = append(embed.Fields, densShield)
	}

	return embed, nil
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
