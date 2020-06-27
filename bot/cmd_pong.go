package bot

import (
    "github.com/bwmarrin/discordgo"
)

func handlePongCmd(s *discordgo.Session, m *discordgo.MessageCreate) error {
    _, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
    return err
}
