package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleInviteCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	embed := b.newEmbed()
	embed.Title = "Rotom-B - Invite Link"
	embed.URL = b.config.Discord.InviteURL

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    "https://i.imgur.com/JqezAUg.gif",
		Width:  500,
		Height: 500,
	}

	embed.Description = fmt.Sprintf(
		"Spread the love of Pokémon by adding Rotom-B to your server with basic message permissions. Enjoy!\n\n"+
			"[Add Rotom-B to your server!](%s)",
		b.config.Discord.InviteURL,
	)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
