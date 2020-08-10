package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// handleHelpCmd handles the "help" command
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

	// Maps on Go do not guarantee key order, so before we fetch help text
	// data, we need to have an alphabetical listing of commands so they
	// print correctly and the same every time.
	commandNames := make([]string, 0)
	for key := range b.commands {
		commandNames = append(commandNames, key)
	}
	sort.Strings(commandNames)

	commandFields := make([]*discordgo.MessageEmbedField, 0)
	for _, name := range commandNames {
		cmd := b.commands[name]
		if cmd.helpText != "" {
			commandFields = append(commandFields, &discordgo.MessageEmbedField{
				Name:   cmd.usage(env.commandPrefix),
				Value:  cmd.helpText,
				Inline: false,
			})
		}
	}

	commandFields = append(commandFields, &discordgo.MessageEmbedField{
		Name: "Support",
		Value: fmt.Sprintf(`[Need more help? Join Rotom-B's support server!](%s)
[If you want to invite Rotom-B to your server, click here](%s)`,
			b.config.Discord.SupportServerURL,
			b.config.Discord.InviteURL,
		),
		Inline: true,
	})

	embed := b.newEmbed()
	embed.Title = "Rotom-B - Help"
	embed.URL = "https://i.imgur.com/xcEamGI.gif"
	embed.Fields = commandFields
	embed.Description = fmt.Sprintf(`
		A list of commands you can use.
		
		* < > Indicate required fields.
		- [ ] Indicate optional fields.
		- Use * for shiny sprites.
		- _Catch Rates are calculated under Raid Specific Conditions: Levels 30-70, 1 HP, and no status modifiers._

		Use "%shelp [command]" for more information about a command.
		`,
		env.commandPrefix,
	)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

// handleCommandUsage sends the embed message with help for a given command
func (b *Bot) handleCommandUsage(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {
	command, ok := b.commands[env.args[0]]
	if !ok {
		_, err := s.ChannelMessageSendEmbed(
			m.ChannelID,
			b.newErrorEmbedf(
				"Command Error",
				`The command "%s" does not exist`,
				env.args[0],
			),
		)
		return err
	}

	aliases := make([]string, 0)
	for key := range b.commands {
		if b.commands[key].alias != "" && b.commands[key].alias == env.args[0] {
			aliases = append(aliases, fmt.Sprintf("%s%s", env.commandPrefix, key))
		}
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Usage",
		Value:  command.usage(env.commandPrefix),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Examples",
		Value:  command.example(env.commandPrefix),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Access",
		Value:  fmt.Sprintf("Admin only: %v", command.adminOnly),
		Inline: true,
	})

	if len(aliases) > 0 {
		aliasesTxt := "`" + strings.Join(aliases, ", ") + "`"
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Aliases",
			Value:  aliasesTxt,
			Inline: true,
		})
	}

	embed := b.newEmbed()
	embed.Title = "Help for **" + env.args[0] + "**"
	embed.URL = "https://i.imgur.com/xcEamGI.gif"
	embed.Description = "**" + env.command + "**: " + command.helpText
	embed.Fields = fields

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
