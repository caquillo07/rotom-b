package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleHelpCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	// if we got a help request for a particular command,
	// return the detailed help details for it instead.
	if len(env.args) > 0 {
		return b.handleCommandUsage(s, env, m)
	}

	embed := &discordgo.MessageEmbed{
		Title: "Rotom-B - Help",
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func (b *Bot) handleCommandUsage(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {
	command, ok := b.commands[env.args[0]]
	if !ok {
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("```The command \"%s\" does not exist```", env.args[0]),
		)
		return err
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Usage",
		Value:  command.usage,
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Examples",
		Value:  command.example,
		Inline: true,
	})

	embed := b.newEmbed()
	embed.Title = "Help for **" + env.args[0] + "**"
	embed.URL = "https://i.imgur.com/xcEamGI.gif"
	embed.Description = "**" + env.command + "**: " + command.helpText
	embed.Fields = fields

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
