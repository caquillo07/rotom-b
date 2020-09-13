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

	// If we are dealing with Chansey or Blissey we need to treat them
	// different, since they are in so many dens they break discord's
	// limit on field size.
	//
	// Both pokemon appear in the same dens for both games, so we will
	// treat both games as one.
	//
	// NOTE: if another pokemon happens to be like this, say in DLC #2,
	// we will need to make this modular
	if pkm.DexID == 113 || pkm.DexID == 242 {
		embed.Fields = getDensEmbedFields("Sword & Shield", swordDensStandard, swordDensHA)
	} else {
		embed.Fields = []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Sword",
				Value:  getDensText(swordDensStandard, swordDensHA),
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Shield",
				Value:  getDensText(shieldDensStandard, shieldDensHA),
				Inline: false,
			},
		}
	}

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

func getDensText(densStandard []string, densHA []string) string {
	var txt string
	for i := 0; i < len(densStandard); i++ {
		den := strings.ToLower(densStandard[i])
		txt += fmt.Sprintf(
			"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml)",
			den,
			den,
		)
		if i != len(densStandard)-1 {
			txt += ", "
		}
	}

	if len(densHA) > 0 {
		txt += "\nHA: "
		for i := 0; i < len(densHA); i++ {
			den := strings.ToLower(densHA[i])
			txt += fmt.Sprintf(
				"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml)",
				den,
				den,
			)
			if i != len(densHA)-1 {
				txt += ", "
			}
		}
	}

	if txt == "" {
		txt = "N/A"
	}
	return txt
}

// TODO(hector): - Make this more DRY when your brain is working again
func getDensEmbedFields(game string, densStandard []string, densHA []string) []*discordgo.MessageEmbedField {
	fields := make([]*discordgo.MessageEmbedField, 0)
	var text string
	const maxDensPerField = 14
	var densInFieldCount int

	// First lets do the standard dens till we run out of space, then we will
	// create a new field and add the rest to that
	for i, den := range densStandard {
		// check to make sure this wont get us out of bounds, if so add it to
		// the fields list, and start over
		isLastInField := densInFieldCount >= maxDensPerField
		if isLastInField {
			name := game
			if len(fields) != 0 {
				name = "(cont'd)"
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   name,
				Value:  text,
				Inline: false,
			})

			text = ""
			densInFieldCount = 0
		}

		densInFieldCount++
		isLastInField = densInFieldCount >= maxDensPerField // update bool flag
		txt := fmt.Sprintf(
			"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml)",
			den,
			den,
		)
		if i != len(densStandard)-1 && !isLastInField {
			txt += ", "
		}
		text += txt
	}

	// now lets create the field, if we added one before due to it being too
	// big, we will add it with a (cont'd) title for contunity
	if text != "" {
		field := &discordgo.MessageEmbedField{
			Name:  game,
			Value: text,
		}
		if len(fields) != 0 {
			field.Name = "(cont'd)"
		}
		fields = append(fields, field)
	}

	// HA is never big enough to fill and entire field (for now), and to keep
	// things from being confusing, lets just add it to its own field.
	if len(densHA) == 0 {
		return fields
	}

	// reset the trackers
	densInFieldCount = 0
	text = ""
	for i, den := range densHA {
		// check to make sure this wont get us out of bounds, if so add it to
		// the fields list, and start over
		isLastInField := densInFieldCount >= maxDensPerField
		if isLastInField {
			name := "HA Dens"
			if len(fields) == 0 {
				name = game + "\n" + name
			}
			if len(fields) != 0 {
				name = "(cont'd)"
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   name,
				Value:  text,
				Inline: false,
			})

			text = ""
			densInFieldCount = 0
		}

		densInFieldCount++
		isLastInField = densInFieldCount >= maxDensPerField // update bool flag
		txt := fmt.Sprintf(
			"[%s](https://www.serebii.net/swordshield/maxraidbattles/den%s.shtml)",
			den,
			den,
		)
		if i != len(densHA)-1 && !isLastInField {
			txt += ", "
		}
		text += txt
	}
	if text != "" {
		field := &discordgo.MessageEmbedField{
			Name:  game,
			Value: text,
		}
		if len(fields) != 0 {
			field.Name = "(cont'd)"
		}
		fields = append(fields, field)
	}
	return fields
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
			"https://caquillo07.github.io/data/dens/den_%s.png",
			strings.ToLower(strings.ReplaceAll(den.Number, " ", "")),
		),
	}
	embed.Fields = []*discordgo.MessageEmbedField{swordField, shieldField}
	return embed, nil
}
