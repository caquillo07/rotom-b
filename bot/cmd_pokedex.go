package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handlePokedexCmd(
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

	pkmArgs := parsePokemonCommand(env.command, env.args)

	// if the name and shininess were not parsed properly, lets assume it
	// follows the order on the help description.
	if pkmArgs.name == "" {
		pkmArgs.name = strings.ReplaceAll(env.args[0], "*", "")
		pkmArgs.isShiny = strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
	}

	pkm, err := b.repository.Pokemon(strings.ToLower(pkmArgs.name))
	if err != nil {
		return botError{
			title:   "Pokémon not found",
			details: fmt.Sprintf("Pokémon %s could not be found.", pkmArgs.name),
		}
	}

	urlPkmName := pkm.Name
	if pkm.DexID == 892 {
		urlPkmName = "Urshifu"
	}

	externalPokedexLinks := fmt.Sprintf(
		"[Bulbapedia Entry](https://bulbapedia.bulbagarden.net/wiki/%s_(Pokémon))\n",
		strings.ReplaceAll(urlPkmName, " ", "_"),
	)

	if len(pkm.Dens.Shield) > 0 || len(pkm.Dens.Sword) > 0 || pkm.Generation == "SwordShield" {
		externalPokedexLinks += fmt.Sprintf(
			"[Serebii Sword & Shield Pokédex](https://serebii.net/pokedex-swsh/%s/)",
			strings.ReplaceAll(strings.ToLower(urlPkmName), " ", ""),
		)
	} else {
		externalPokedexLinks += fmt.Sprintf(
			"[Serebii Sun & Moon Pokédex](https://serebii.net/pokedex-sm/%03d.shtml)",
			pkm.DexID,
		)
	}

	abilities := "`" + pkm.Abilities.Ability1
	if pkm.Abilities.Ability2 != "" {
		abilities += ",\n" + pkm.Abilities.Ability2
	}
	if pkm.Abilities.AbilityH != "" {
		abilities += ",\n" + pkm.Abilities.AbilityH + " (HA)"
	}
	abilities += "`"

	eggGroups := pkm.EggGroup1
	if pkm.EggGroup2 != "" {
		eggGroups = fmt.Sprintf(
			"%s, %s",
			pkm.EggGroup1,
			pkm.EggGroup2,
		)
	}

	forms := createJoinedPkmInfo("Forms", pkm.Forms)
	densSword := createJoinedPkmInfo("Sword", pkm.Dens.Sword)
	densShield := createJoinedPkmInfo("Shield", pkm.Dens.Shield)

	embed := b.newEmbed()
	embed.Title = fmt.Sprintf("%s Pokédex Info", strings.Title(pkmArgs.name))
	embed.Image = &discordgo.MessageEmbedImage{
		URL:    pkm.SpriteImage(pkmArgs.isShiny, pkmArgs.form),
		Width:  300,
		Height: 300,
	}

	embed.URL = fmt.Sprintf(
		"https://bulbapedia.bulbagarden.net/wiki/%s_(Pokémon)",
		strings.ReplaceAll(urlPkmName, " ", "_"),
	)

	types := pkm.Type1
	if pkm.Type2 != "" {
		types += " / " + pkm.Type2
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name: "Base Stats",
			Value: fmt.Sprintf(
				"HP: `%d`\n"+
					"Atk: `%d`\n"+
					"Def: `%d`\n"+
					"Spa: `%d`\n"+
					"SpD: `%d`\n"+
					"Spe: `%d`\n"+
					"Total: `%d`",
				pkm.BaseStats.HP,
				pkm.BaseStats.Atk,
				pkm.BaseStats.Def,
				pkm.BaseStats.SpA,
				pkm.BaseStats.SpD,
				pkm.BaseStats.Spd,
				pkm.BaseStats.Total,
			),
			Inline: true,
		},
		{
			Name:   "Abilities",
			Value:  abilities,
			Inline: true,
		},
		{
			Name: "Pokémon Misc. Info",
			Value: fmt.Sprintf(
				"Type: `%s`\n"+
					"Gender Ratio: `%s`\n"+
					"Height / Weight: `%sm / %skg`\n"+
					"Catch Rate: `%d`\n"+
					"Generation: `%s`\n"+
					"Egg Groups: `%s`\n"+
					"%s",
				types,
				pkm.GenderRatio,
				fmt.Sprintf("%.2f", pkm.Height),
				fmt.Sprintf("%.2f", pkm.Weight),
				pkm.CatchRate,
				pkm.Generation,
				eggGroups,
				forms,
			),
			Inline: true,
		},
	}

	if len(pkm.Dens.Shield) > 0 || len(pkm.Dens.Sword) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Dens",
			Value: fmt.Sprintf(
				"%s\n%s",
				densSword,
				densShield,
			),
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "More Info",
		Value:  externalPokedexLinks,
		Inline: false,
	})

	return sendEmbed(s, m.ChannelID, embed)
}

func createJoinedPkmInfo(prefix string, info []string) string {
	joinedInfo := ""
	if len(info) > 0 {
		joinedInfo = prefix + ": `" + strings.Join(info, ", ") + "`"
	}
	return joinedInfo
}
