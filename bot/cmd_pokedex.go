package bot

import (
	"fmt"
	"strconv"
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

	var embed *discordgo.MessageEmbed
	var err error

	pkmnName := env.args[0]

	pkmn, err := b.pokemonRepo.pokemon(strings.ToLower(pkmnName))

	if err != nil {
		return botError{
			title: "Pokémon not found",
			details: fmt.Sprintf("Pokémon %s could not be found.",
				pkmnName),
		}
	}

	var forms string
	if len(pkmn.Forms) > 0 {
		forms = "Forms: " + strings.Join(pkmn.Forms, ", ") + "."
	} else {
		forms = ""
	}

	abilities := pkmn.Abilities.Ability1
	if pkmn.Abilities.Ability2 != "" {
		abilities += ", " + pkmn.Abilities.Ability2
	}

	externalPokedexLinks := fmt.Sprintf(
		"[Bulbapedia Entry](https://bulbapedia.bulbagarden.net/wiki/%s_(Pokémon))\n",
		strings.Title(pkmn.Name),
	)
	if len(pkmn.Dens.Shield) > 0 || len(pkmn.Dens.Sword) > 0 || pkmn.Generation == "SwordShield" {
		externalPokedexLinks += fmt.Sprintf(
			"[Serebii Sword & Shield Pokédex](https://serebii.net/pokedex-swsh/%s/)",
			strings.ToLower(pkmn.Name),
		)
	} else {
		externalPokedexLinks += fmt.Sprintf(
			"[Serebii Sun & Moon Pokédex](https://serebii.net/pokedex-sm/%s.shtml)",
			strconv.Itoa(pkmn.DexID),
		)
	}

	embed = b.newEmbed()
	embed.Title = fmt.Sprintf("%s Pokédex Info", strings.Title(pkmnName))
	embed.Description = fmt.Sprintf("%s Pokédex info", strings.Title(pkmn.Name))

	embed.Image = &discordgo.MessageEmbedImage{
		URL:    pkmn.spriteImage(false, ""),
		Width:  300,
		Height: 300,
	}

	embed.URL = fmt.Sprintf(
		"https://bulbapedia.bulbagarden.net/wiki/%s_(Pokémon)",
		strings.ToLower(pkmn.Name),
	)

	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name: "Pokémon Mis. Info",
			Value: fmt.Sprintf(
				"Gender Ratio: `%s`\n"+
					"Height / Weight: `%s / %s`\n"+
					"Catch Rate: `%s`\n"+
					"Generation: `%s`\n"+
					"Egg Groups: `%s, %s`\n"+
					"`%s`",
				pkmn.GenderRatio,
				fmt.Sprintf("%.2f", pkmn.Height),
				fmt.Sprintf("%.2f", pkmn.Weight),
				strconv.Itoa(pkmn.CatchRate),
				pkmn.Generation,
				pkmn.EggGroup1,
				pkmn.EggGroup2,
				forms,
			),
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name: "Base Stats",
			Value: fmt.Sprintf(
				"HP: `%s`\n"+
					"Atk: `%s`\n"+
					"Def:  `%s`\n"+
					"Spa: `%s`\n"+
					"SpD: `%s`\n"+
					"Spe: `%s`\n",
				strconv.Itoa(pkmn.BaseStats.HP),
				strconv.Itoa(pkmn.BaseStats.Atk),
				strconv.Itoa(pkmn.BaseStats.Def),
				strconv.Itoa(pkmn.BaseStats.SpA),
				strconv.Itoa(pkmn.BaseStats.SpD),
				strconv.Itoa(pkmn.BaseStats.Spd),
			),
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name: "Abilities",
			Value: fmt.Sprintf(
				"Abilities: `%s`\n"+
					"Hidden Ability: `%s`",
				abilities,
				pkmn.Abilities.AbilityH,
			),
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name: "Dens",
			Value: fmt.Sprintf(
				"Sword: `%s`\n"+
					"Shield: `%s`",
				strings.Join(pkmn.Dens.Sword, ", "),
				strings.Join(pkmn.Dens.Sword, ", "),
			),
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "See More Information",
			Value:  externalPokedexLinks,
			Inline: true,
		},
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
