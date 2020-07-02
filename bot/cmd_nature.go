package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleNatureCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	natures := make(map[string]string)

	natures["Hardy"] = "No changes"
	natures["Lonely"] = "+Atk -Def"
	natures["Brave"] = "+Atk -Spe"
	natures["Adamant"] = "+Atk -SpA"
	natures["Naughty"] = "+Atk -SpD"
	natures["Bold"] = "+Def -Atk"
	natures["Docile"] = "No changes"
	natures["Relaxed"] = "+Def -Spe"
	natures["Impish"] = "+Def -Atk"
	natures["Lax"] = "+Def -SpD"
	natures["Timid"] = "+Spe -Atk"
	natures["Hasty"] = "+Spe -Def"
	natures["Serious"] = "No changes"
	natures["Jolly"] = "+Spe -SpA"
	natures["Naive"] = "+Spe -SpD"
	natures["Modest"] = "+SpA -Atk"
	natures["Mild"] = "+SpA -Def"
	natures["Quiet"] = "+SpA -Spe"
	natures["Bashful"] = "No changes"
	natures["Rash"] = "+SpA -SpD"
	natures["Calm"] = "+SpD -Atk"
	natures["Gentle"] = "*+SpD -Def"
	natures["Sassy"] = "+SpD -Spe"
	natures["Careful"] = "+SpD -SpA"
	natures["Quirky"] = "No changes"

	var embed *discordgo.MessageEmbed
	var err error

	embed = b.newEmbed()

	if len(env.args) == 0 {
		embed.Title = "Pokémon Natures Chart (from Bulbapedia)"
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://raphgg.github.io/den-bot/data/icons/natures.PNG",
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	}

	nature := env.args[0]

	if natureInfo, ok := natures[strings.Title(nature)]; ok {
		embed.Title = fmt.Sprintf("%s Nature Info", strings.Title(nature))
		embed.Description = natureInfo
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err

	} else {
		return botError{
			title: "Nature name not found",
			details: fmt.Sprintf("Nature %s could not be found.",
				nature),
		}
	}
}
