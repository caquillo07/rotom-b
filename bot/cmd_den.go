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

	var embed *discordgo.MessageEmbed
	var err error

	_, err = strconv.Atoi(env.args[0])
	isNumber := err == nil

	if isNumber {
		embed, err = b.getDenFromNumber(env.args[0])
	} else {
		embed, err = b.getDensFromPokemon(env.args[0])
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func (b *Bot) getDensFromPokemon(pkmnName string) (*discordgo.MessageEmbed, error) {

	pokemon, err := b.pokemonRepo.pokemon(strings.ToLower(pkmnName))
	if err != nil {
		return nil, err
	}

	embed := b.newEmbed()
	embed.Title = pokemon.Name + " is in the following Dens:"
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf(
			"https://raphgg.github.io/den-bot/data/sprites/pokemon/normal/%s.gif",
			strings.ToLower(strings.ReplaceAll(pokemon.Name, " ", "")),
		),
	}

	swordField := &discordgo.MessageEmbedField{}
	swordField.Inline = false
	swordField.Name += "*Sword:* "
	for i := 0; i < len(pokemon.Dens.Sword); i++ {
		if i == len(pokemon.Dens.Sword)-1 {
			swordField.Value += fmt.Sprintf(
				"[%s](https://raphgg.github.io/den-bot/data/sprites/pokemon/normal/%s.gif)",
				strings.ToLower(pokemon.Dens.Sword[i]),
				strings.ToLower(strings.ReplaceAll(pokemon.Name, " ", "")),
			)
			break
		}
		swordField.Value += fmt.Sprintf(
			"[%s](https://raphgg.github.io/den-bot/data/sprites/pokemon/normal/%s.gif), ",
			strings.ToLower(pokemon.Dens.Sword[i]),
			strings.ToLower(strings.ReplaceAll(pokemon.Name, " ", "")),
		)
	}

	shieldField := &discordgo.MessageEmbedField{}
	shieldField.Inline = false
	shieldField.Name += "*Shield:* "
	for i := 0; i < len(pokemon.Dens.Shield); i++ {
		if i == len(pokemon.Dens.Shield)-1 {
			shieldField.Value += fmt.Sprintf(
				"[%s](https://raphgg.github.io/den-bot/data/sprites/pokemon/normal/%s.gif)",
				strings.ToLower(pokemon.Dens.Shield[i]),
				strings.ToLower(strings.ReplaceAll(pokemon.Name, " ", "")),
			)
			break
		}
		shieldField.Value += fmt.Sprintf(
			"[%s](https://raphgg.github.io/den-bot/data/sprites/pokemon/normal/%s.gif), ",
			strings.ToLower(pokemon.Dens.Shield[i]),
			strings.ToLower(strings.ReplaceAll(pkmnName, " ", "")),
		)
	}
	embed.Fields = []*discordgo.MessageEmbedField{swordField, shieldField}

	return embed, err
}

func (b *Bot) getDenFromNumber(denNumber string) (*discordgo.MessageEmbed, error) {

	den, err := b.pokemonRepo.den(denNumber)
	if err != nil {
		return nil, err
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
	// msgEmbed.Description = message

	return embed, err
}
