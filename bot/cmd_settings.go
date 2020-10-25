package bot

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/caquillo07/rotom-bot/repository"
)

const (
	timeConverterURL = "https://www.worldtimeserver.com/convert_time_in_UTC.aspx"

	actionAdd    = "add"
	actionRemove = "remove"
	actionReset  = "reset"
)

var channelIDRegex = regexp.MustCompile("<#(\\w+)>")

func (b *Bot) handleConfigCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return err
	}

	guildSettings, err := b.getOrCreateGuildSettings(guild)
	if err != nil {
		return err
	}

	if len(env.args) == 0 {
		embed, err := b.currentSettingsEmbed(s, guildSettings)
		if err != nil {
			return err
		}
		return sendEmbed(s, m.ChannelID, embed)
	}

	switch c := env.args[0]; c {
	case "prefix":
		if len(env.args) < 2 || env.args[1] == "" {
			return botError{
				title:   "Validation Error",
				details: "Prefix is required to update the setting",
			}
		}
		guildSettings.BotPrefix = env.args[1]
	case "listen":
		if len(env.args) < 2 || env.args[1] == "" {
			return botError{
				title:   "Validation Error",
				details: "Channel or action is required to update the setting",
			}
		}

		// ignored the bad ones, but we can return this as a warning later.
		_, err := handleListenUpdate(s, env.args, guildSettings)
		if err != nil {
			return err
		}

	default:
		return botError{
			title:   "Validation Error",
			details: c + " is not a valid setting",
		}
	}
	guildSettings.LastUpdatedBy = m.Author.ID
	if err := b.repository.UpdateGuildSettings(guildSettings); err != nil {
		return err
	}

	embed := b.newEmbed()
	embed.Title = "Update Successful"
	embed.Description = "Setting was updated successfully"
	embed.Color = 0x00FF00
	return sendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) currentSettingsEmbed(s *discordgo.Session, settings *repository.GuildSettings) (*discordgo.MessageEmbed, error) {
	lastUpdatedBy := "Never"
	if settings.LastUpdatedBy != "" {
		user, err := s.User(settings.LastUpdatedBy)
		if err != nil {
			return nil, err
		}
		lastUpdatedBy = user.String()

		if t := settings.UpdatedAt; !t.IsZero() {
			lastUpdatedBy += fmt.Sprintf(
				" on [%s](%s?y=%d&mo=%d&d=%d&h=%d&mn=%d)",
				t.Format(time.RFC822),
				timeConverterURL,
				t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
			)
		}
	}

	embed := b.newEmbed()
	embed.Title = settings.Name + "'s settings"
	embed.Description = "List of settings and their values for server"
	listeningOn := ""
	for i, c := range settings.ListeningChannels {
		listeningOn += c.Name
		if len(settings.ListeningChannels)-1 != i {
			listeningOn += ", "
		}
	}
	if listeningOn == "" {
		listeningOn = "N/A"
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Prefix",
			Value:  settings.BotPrefix,
			Inline: false,
		},
		{
			Name:   "Listening on",
			Value:  listeningOn,
			Inline: false,
		},
		{
			Name:   "Last Updated By",
			Value:  lastUpdatedBy,
			Inline: false,
		},
	}

	return embed, nil
}

func getActionFromArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	c := args[1]
	if c != actionAdd && c != actionRemove && c != actionReset {
		return ""
	}

	return c
}

func handleListenUpdate(s *discordgo.Session, args []string, settings *repository.GuildSettings) ([]string, error) {
	action := getActionFromArgs(args)
	// if we are resetting, just empty and return early
	if action == actionReset {
		settings.ListeningChannels = nil
		return nil, nil
	}

	// if we get an action but no channels
	if action != "" && len(args) < 3 {
		return nil, botError{
			title:   "Validation Error",
			details: "Channel(s) are required to update the setting",
		}
	}

	// If there is no action first, or there is only one channel and its
	// malformed, return an error
	idx := 1
	if action != "" {
		idx = 2
	}
	channelIDs := args[idx:]
	cleanedIDs := make([]string, 0, len(channelIDs))
	badChannelIDs := make([]string, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		if match := channelIDRegex.FindStringSubmatch(channelID); len(match) == 2 && match[1] != "" {
			cleanedIDs = append(cleanedIDs, match[1])
			continue
		}
		badChannelIDs = append(badChannelIDs, channelID)
	}

	if len(cleanedIDs) == 0 && len(channelIDs) == 1 {
		return nil, botError{
			title:   "Validation Error",
			details: fmt.Sprintf("Channel %q is not a valid channel", args[1]),
		}
	}

	finalChannels := make([]*repository.GuildSettingChannel, len(cleanedIDs))
	for i, id := range cleanedIDs {
		channel, err := s.Channel(id)
		if err != nil {
			return nil, err
		}
		finalChannels[i] = &repository.GuildSettingChannel{
			ID:   channel.ID,
			Name: channel.Name,
		}
	}

	// if we have no actions, but cleaned IDs, just do a replace all
	if action == "" {
		settings.ListeningChannels = finalChannels
		return badChannelIDs, nil
	}

	// At this point, just run through the list and perform the desired action
	for _, channel := range finalChannels {
		switch action {
		case actionAdd:
			settings.ListeningChannels = append(settings.ListeningChannels, channel)
		case actionRemove:
			cleanedList := make([]*repository.GuildSettingChannel, 0, len(settings.ListeningChannels))
			for _, c := range settings.ListeningChannels {
				if c.ID == channel.ID {
					continue
				}
				cleanedList = append(cleanedList, c)
			}
			settings.ListeningChannels = cleanedList
		case "":
			// ignore.
		}
	}

	return badChannelIDs, nil
}
