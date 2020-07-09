package bot

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handlePokedexCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	if len(env.args) == 0 {
		return botError{
			title:   "Validation Error",
			details: "Please enter a den number or a Pokémon name to look for related dens.",
		}
	}

	var err error

	return err

}
