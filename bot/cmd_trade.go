// added by Lewymd, if it breaks it's his fault

package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// handleHelpCmd handles the "trade" command

func (b *Bot) handleTradeCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	embed := b.newEmbed()
	embed.Title = "Rotom-B - trade"
	embed.URL = "https://i.imgur.com/NQFaIYq.png"
	embed.Description = fmt.Sprintf(`
		A list of universal trade codes to help you find trade partners
		
		Ditto - 4448 4448
		Version Exclusives - 2211 2211
		Breedject Trades - 1337 1337
		Shiny Trades - 0000 0777
		Trade requirement Evolutions - 1234 4321
		`,
	)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}