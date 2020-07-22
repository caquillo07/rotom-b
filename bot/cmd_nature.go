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

	natures := map[string]string{
		"Hardy":   "No changes", // nolint:goconst
		"Lonely":  "+Atk -Def",
		"Brave":   "+Atk -Spe",
		"Adamant": "+Atk -SpA",
		"Naughty": "+Atk -SpD",
		"Bold":    "+Def -Atk",
		"Docile":  "No changes",
		"Relaxed": "+Def -Spe",
		"Impish":  "+Def -SpA",
		"Lax":     "+Def -SpD",
		"Timid":   "+Spe -Atk",
		"Hasty":   "+Spe -Def",
		"Serious": "No changes",
		"Jolly":   "+Spe -SpA",
		"Naive":   "+Spe -SpD",
		"Modest":  "+SpA -Atk",
		"Mild":    "+SpA -Def",
		"Quiet":   "+SpA -Spe",
		"Bashful": "No changes",
		"Rash":    "+SpA -SpD",
		"Calm":    "+SpD -Atk",
		"Gentle":  "+SpD -Def",
		"Sassy":   "+SpD -Spe",
		"Careful": "+SpD -SpA",
		"Quirky":  "No changes",
	}

	embed := b.newEmbed()

	if len(env.args) == 0 {
		embed.Title = "Pokémon Natures Chart (from Bulbapedia)"
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://raphgg.github.io/den-bot/data/icons/natures.PNG",
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	}

	nature := env.args[0]
	if natureInfo, ok := natures[strings.Title(nature)]; ok {
		embed.Title = fmt.Sprintf("%s Nature Info", strings.Title(nature))
		embed.Description = natureInfo
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	}

	return botError{
		title:   "Nature name not found",
		details: fmt.Sprintf("Nature %s could not be found.", nature),
	}
}
