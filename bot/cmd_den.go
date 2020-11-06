package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/caquillo07/rotom-bot/repository"
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

	pkmArgs := parsePokemonCommand(env.command, env.args)

	if pkmArgs.den != "" {
		embed, err := b.getDenFromNumber(env.args[0])
		if err != nil {
			return err
		}
		return sendEmbed(s, m.ChannelID, embed)
	}

	// if the name and shininess were not parsed properly, lets assume it
	// follows the order on the help description.
	if pkmArgs.name == "" {
		pkmArgs.name = strings.ReplaceAll(env.args[0], "*", "")
		pkmArgs.isShiny = strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
	}

	embed, err := b.getDensFromPokemon(pkmArgs.name, pkmArgs.form, pkmArgs.isShiny)
	if err != nil {
		return err
	}

	return sendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) getDensFromPokemon(pkmnName, form string, isShiny bool) (*discordgo.MessageEmbed, error) {

	// the pokemon.json file has the galarian pokemon as a separate pokemon,
	// so in this special case we must append it to the name. Maybe we need to
	// change this?
	if form == galarian {
		pkmnName = form + " " + pkmnName
	}
	pkm, err := b.repository.Pokemon(pkmnName)
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

	swordDensHA := make([]string, 0)
	swordDensStandard := make([]string, 0)
	shieldDensHA := make([]string, 0)
	shieldDensStandard := make([]string, 0)

	// Sword
	for _, d := range pkm.Dens.Sword {
		den, err := b.repository.Den(d)
		if err != nil {
			return nil, nil
		}

		if isDenPokemonHA(pkm.Name, den.Sword) {
			swordDensHA = append(swordDensHA, d)
		} else {
			swordDensStandard = append(swordDensStandard, d)
		}
	}

	// Shield
	for _, d := range pkm.Dens.Shield {
		den, err := b.repository.Den(d)
		if err != nil {
			return nil, nil
		}

		if isDenPokemonHA(pkm.Name, den.Shield) {
			shieldDensHA = append(shieldDensHA, d)
		} else {
			shieldDensStandard = append(shieldDensStandard, d)
		}
	}

	embed := b.newEmbed()
	embed.Title = pkm.Name + " is in the following Dens:"
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    pkm.SpriteImage(isShiny, form),
		Width:  150,
		Height: 150,
	}

	embed.Fields = append(embed.Fields, getDenFields("Sword", swordDensStandard)...)
	embed.Fields = append(embed.Fields, getDenFields("Sword HA", swordDensHA)...)
	embed.Fields = append(embed.Fields, getDenFields("Shield", shieldDensStandard)...)
	embed.Fields = append(embed.Fields, getDenFields("Shield HA", shieldDensHA)...)
	return embed, nil
}

func isDenPokemonHA(pkmName string, gameDens []*repository.DenPokemon) bool {
	for _, gd := range gameDens {
		if pkmName != gd.Name || !strings.HasPrefix(gd.Ability, "Hidden") {
			continue
		}
		return true
	}
	return false
}

func getDenFields(title string, dens []string) []*discordgo.MessageEmbedField {
	getField := func(value string, isExtra bool) *discordgo.MessageEmbedField {
		name := title
		if isExtra {
			name += " (cont'd)"
		}
		return &discordgo.MessageEmbedField{
			Name:   name,
			Value:  strings.TrimSuffix(value, ", "),
			Inline: false,
		}
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	var text string
	for i := 0; i < len(dens); i++ {
		den := fmt.Sprintf(
			// its ok to add the comma at the end here, it will get trimmed by
			// the field create function
			"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml), ",
			dens[i],
			dens[i],
		)

		// if we have reached the maximum allowed characters, its time to make
		// a new field.
		if len(text)+len(den) >= embedFieldValueMaxCharacters {
			fields = append(fields, getField(text, len(fields) != 0))
			text = ""
		}

		text += den
	}
	if text == "" {
		text = "N/A"
	}

	// we need to make sure we add any leftovers, or the first field if we
	// never reached the max characters, or if its just N/A
	return append(fields, getField(text, len(fields) != 0))
}

func (b *Bot) getDenFromNumber(denNumber string) (*discordgo.MessageEmbed, error) {

	den, err := b.repository.Den(denNumber)
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
	if swordField.Value == "" {
		swordField.Value = "N/A"
	}

	shieldField := &discordgo.MessageEmbedField{}
	shieldField.Inline = true
	shieldField.Name += "HA in Shield"
	for i := 0; i < len(den.Shield); i++ {
		if den.Shield[i].Ability != "Standard" {
			shieldField.Value += den.Shield[i].Name + "\n"
		}
	}
	if shieldField.Value == "" {
		shieldField.Value = "N/A"
	}

	embed := b.newEmbed()
	embed.Title = "Pokémon found in Den " + den.Number + ":"
	embed.URL = fmt.Sprintf(
		"https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml",
		strings.ToLower(strings.ReplaceAll(den.Number, " ", "")),
	)
	embed.Image = &discordgo.MessageEmbedImage{
		URL: fmt.Sprintf(
			"https://raw.githubusercontent.com/caquillo07/rotom-b-data/master/dens/den_%s.png",
			strings.ToLower(strings.ReplaceAll(den.Number, " ", "")),
		),
	}
	embed.Fields = []*discordgo.MessageEmbedField{swordField, shieldField}
	return embed, nil
}
