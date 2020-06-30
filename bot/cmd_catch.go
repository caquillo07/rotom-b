package bot

import (
	"github.com/bwmarrin/discordgo"
)

// handleCatchCmd handles the catch command, sends back a
// detailed summary of catch rates for a given Pokémon & Ball.
func (b *Bot) handleCatchCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {
	if len(env.args) == 0 {
		return botError{
			title:   "Validation Error",
			details: "Please enter a Pokémon to catch followed by a Poké-Ball of your choice.",
		}
	}

	return nil
}
