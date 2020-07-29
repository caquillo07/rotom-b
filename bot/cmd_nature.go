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
	natureInfo, ok := natures[strings.Title(nature)]
	if !ok {
		return botError{
			title:   "Nature name not found",
			details: fmt.Sprintf("Nature %s could not be found.", nature),
		}
	}

	embed.Title = fmt.Sprintf("%s Nature Info", strings.Title(nature))
	embed.Description = natureInfo
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
