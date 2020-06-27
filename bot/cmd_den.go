package bot

import (
    "github.com/bwmarrin/discordgo"
)

func (b *Bot) handleDenCmd(s *discordgo.Session, m *discordgo.MessageCreate) error {
    _, err := s.ChannelMessageSend(m.ChannelID, "Den command!")
    return err
}
