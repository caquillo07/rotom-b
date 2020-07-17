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

	inviteUrl := "https://discordapp.com/oauth2/authorize?client_id=726478988276531212&scope=bot&permissions=52224"

	embed := b.newEmbed()
	embed.Title = "Rotom-B - Invite Link"
	embed.URL = inviteUrl

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:    "https://i.imgur.com/JqezAUg.gif",
		Width:  500,
		Height: 500,
	}

	embed.Description = fmt.Sprintf(
		"Spread the love of Pokémon by adding Rotom-B to your server with basic message permissions. Enjoy!\n\n"+
			"[Add Rotom-B to your server!](%s)",
		inviteUrl,
	)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
