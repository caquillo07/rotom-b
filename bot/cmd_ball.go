package bot

import (
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// handleBallCmd handles the ball command, sends back a gif
// animation of the ball being used, and basic information.
func (b *Bot) handleBallCmd(
	s *discordgo.Session,
	env *commandEnvironment,
	m *discordgo.Message,
) error {

	ball, err := b.pokemonRepo.ball(env.args[0])
	if err != nil {
		return err
	}

	embed := b.newEmbed()
	embed.Title = fmt.Sprintf("**__%s__**", ball.Name)
	embed.URL = fmt.Sprintf(
		"https://serebii.net/itemdex/%s.shtml",
		strings.ToLower(strings.ReplaceAll(ball.Name, " ", "")),
	)
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: fmt.Sprintf(
			"https://raphgg.github.io/den-bot/data/sprites/balls/%s.png",
			ball.ID,
		),
	}
	embed.Image = &discordgo.MessageEmbedImage{
		URL: fmt.Sprintf(
			"https://raphgg.github.io/den-bot/data/sprites/balls/%s.gif?cache=42069",
			ball.ID,
		),
	}

	// If the formatter is a whole number, we will want to print it as such.
	// If its a decimal value, we will print it with one decimal precision.
	modFormatter := "%.0f"
	if ball.Modifier != math.Trunc(ball.Modifier) {
		modFormatter = "%.1f"
	}
	embed.Description = fmt.Sprintf(`
		**Ball Effects:** %s
		**Ball Modifier:** %sx
		**Ball Conditions:** %s
		`,
		ball.Effect,
		fmt.Sprintf(modFormatter, ball.Modifier), // this is confusing, but it works
		ball.Conditions,
	)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
