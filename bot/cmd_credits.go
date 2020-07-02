package bot

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleCreditsCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	embed := b.newEmbed()
	embed.Title = "Rotom-B - Credits"
	embed.URL = "https://github.com/caquillo07/rotom-b"
	embed.Description = "Thanks to everyone who helped with the bot! Without them this would not be possible:"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "Bot Development",
			Value:  "This bot was developed by [@caquillo07](https://github.com/caquillo07/) and [@hypermilla](https://twitter.com/hypermilla). It's open source and you can [check it out on GitHub](https://github.com/caquillo07/rotom-b).",
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "External Resources",
			Value:  "[Serebii](https://serebii.net/) & [Bulbapedia](https://bulbapedia.bulbagarden.net/) for their amazing base of Pokémon information and all their years of hosting it for the community.",
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Sprite Work",
			Value:  "[PkParaíso](https://pkparaiso.com/) & [Ian Clail](https://www.smogon.com/forums/threads/sun-moon-sprite-project.3577711/) (Layell) for the awesome sprite work on 1300+ Pokémon and different forms. Also, Tax for the ball animation sprites. These sprites were taken from the Alcremie-B Den Bot and re-used.",
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Den Bot Inspiration and Original Concept",
			Value:  "[Alcremie-B Bot by Droopy](https://github.com/RaphGG/den-bot) is the inspiration behind this bot, which would not have been possible without the original open-source work.",
			Inline: false,
		},
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
