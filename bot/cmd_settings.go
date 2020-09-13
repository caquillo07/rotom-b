package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/caquillo07/rotom-bot/repository"
)

const timeConverterURL = "https://www.worldtimeserver.com/convert_time_in_UTC.aspx"

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

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Prefix",
			Value:  settings.BotPrefix,
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
