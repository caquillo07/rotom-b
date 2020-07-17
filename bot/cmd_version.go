package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/caquillo07/rotom-bot/metrics"
)

func (b *Bot) handleVersionCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	url := "https://github.com/caquillo07/rotom-b"

	embed := b.newEmbed()
	embed.Title = "Rotom-B - Bot Version"
	embed.URL = url

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    "https://i.imgur.com/Xc06rex.jpg",
		Width:  500,
		Height: 500,
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Version: %s", metrics.Version),
			Inline: false,
			Value: fmt.Sprintf(
				"Running version `%s` of Rotom-B. To get the latest version, [head to our github page.](%s)",
				metrics.Version,
				url,
			),
		},
		&discordgo.MessageEmbedField{
			Name:   "Commit",
			Inline: false,
			Value: fmt.Sprintf(
				"Running commit `%s`\n"+
					"on Branch `%s`",
				metrics.Commit,
				metrics.Branch,
			),
		},
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
