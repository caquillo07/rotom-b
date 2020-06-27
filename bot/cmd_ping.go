package bot

import (
    "github.com/bwmarrin/discordgo"
)

func (b *Bot) handlePingCmd(s *discordgo.Session, m *discordgo.MessageCreate) error {
    _, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
    return err
}
