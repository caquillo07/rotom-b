//todo requested by Hector - add Lewymd to contributors list
//added by Lewymd, if it breaks it's his fault

package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// handleTradeCmd handles the "trade" command

func (b *Bot) handleTradeCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	embed := b.newEmbed()
	embed.Title = "Rotom-B - Trade"
	embed.URL = "https://i.imgur.com/NQFaIYq.png"
	
//this needs to be tidied up and probably use a function to check if values exist then produce a string based on that however even in the current form won't print blank lines but is hard limited to 10 inputs
	embed.Description = fmt.Sprintf(`
		A list of universal trade codes to help you find trade partners

		`,
	) + b.config.TradeCodes.TradeDesc + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc2 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code2 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc3 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code3 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc4 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code4 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc5 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code5 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc6 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code6 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc7 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code7 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc8 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code8 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc9 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code9 + fmt.Sprintf(`
	`,) + b.config.TradeCodes.TradeDesc10 + fmt.Sprintf(`	`,) +  b.config.TradeCodes.Code10 + fmt.Sprintf(`
	`,)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}