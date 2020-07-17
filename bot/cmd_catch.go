package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	catchRateConfidenceURL = "https://bulbapedia.bulbagarden.net/wiki/Catch_rate#Probability_of_capture"
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

	isShiny := strings.HasSuffix(env.args[0], "*") || strings.HasPrefix(env.args[0], "*")
	cleanPkmName := strings.ReplaceAll(env.args[0], "*", "")
	pkm, err := b.pokemonRepo.pokemon(cleanPkmName)
	if err != nil {
		return botError{
			title:   "Not Found",
			details: fmt.Sprintf("Pokemon %s was not found", cleanPkmName),
		}
	}

	form := getFormFromArgs(env.args)
	ball, err := b.pokemonRepo.ball(env.args[len(env.args)-1])
	if err != nil && err != errBallDoesNotExist {
		return err
	}

	// now make sure that someone did not just send in a ball and no pokemon
	if ball != nil && len(env.args) == 1 {
		return botError{
			title:   "Validation Error",
			details: "A pokemon must be provided alongside the PokeBall for the catch command",
		}
	}

	// If the ball does not exist, that means we got just a pokemon request
	if ball == nil {
		embed, err := b.getPokemonTopFourBalls(pkm, form, isShiny)
		if err != nil {
			return err
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	}

	// If we got a ball, we are doing an specific check against a pokemon.
	embed, err := b.getPokemonCatchRate(pkm, ball, form, isShiny)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func (b *Bot) getPokemonCatchRate(
	pkm *pokemon,
	ball *pokeBall,
	form string,
	shiny bool,
) (*discordgo.MessageEmbed, error) {
	form = getSpriteForm(form)
	name := pkm.Name
	isGmax := form == gigantamax
	if isGmax {
		name = "G-Max " + name
	}

	embed := b.newEmbed()
	embed.Title = fmt.Sprintf("%s catch probability", name)
	embed.Color = b.getPokemonColor(pkm.Type1)
	embed.Image = &discordgo.MessageEmbedImage{
		URL: pkm.spriteImage(shiny, form),
	}

	ball.Modifier = ball.catchModifier(pkm)
	lowerCatchProb, lowerConfidence := pkm.captureRate(ball, 30, 0, isGmax, false)
	higherCatchProb, highConfidence := pkm.captureRate(ball, 70, 31, isGmax, false)
	description := fmt.Sprintf("%.2f%%", lowerCatchProb)
	if lowerCatchProb != higherCatchProb {
		description += fmt.Sprintf(" ~ %.2f%%", higherCatchProb)
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name: ball.Name,
			Value: fmt.Sprintf(
				"%s catch rate: `%s`\n[Confidence level](%s): %s",
				ball.Name,
				description,
				catchRateConfidenceURL,
				confidenceEmoji(lowerConfidence && highConfidence),
			),
			Inline: true,
		},
	}

	return embed, nil
}

func (b *Bot) getPokemonTopFourBalls(pkm *pokemon, form string, shiny bool) (*discordgo.MessageEmbed, error) {
	form = getSpriteForm(form)
	name := pkm.Name
	isGmax := form == gigantamax
	if isGmax {
		name = "G-Max " + name
	}

	embed := b.newEmbed()
	embed.Title = "Best Catch Rates"
	embed.Description = fmt.Sprintf("The best balls for catching %s are:", name)
	embed.Color = b.getPokemonColor(pkm.Type1)
	embed.Image = &discordgo.MessageEmbedImage{
		URL: pkm.spriteImage(shiny, form),
	}

	topFourDesp := make([]string, 0)
	for _, pb := range b.getBestBallsForPokemon(pkm, isGmax) {
		lowerCatchProb, _ := pkm.captureRate(pb, 30, 0, isGmax, false)
		higherCatchProb, _ := pkm.captureRate(pb, 70, 31, isGmax, false)
		description := fmt.Sprintf(
			"%s: `%.2f%%",
			pb.Name,
			lowerCatchProb,
		)
		if lowerCatchProb != higherCatchProb {
			description += fmt.Sprintf(" ~ %.2f%%", higherCatchProb)
		}
		description += "`"
		topFourDesp = append(topFourDesp, description)
	}

	pkBall, err := b.pokemonRepo.ball("poke ball")
	if err != nil {
		return nil, err
	}

	stdLowerCatchProb, _ := pkm.captureRate(pkBall, 30, 0, isGmax, false)
	stdHigherCatchProb, _ := pkm.captureRate(pkBall, 70, 31, isGmax, false)
	standardDescription := fmt.Sprintf("%.2f%%", stdLowerCatchProb)
	if stdLowerCatchProb != stdHigherCatchProb {
		standardDescription += fmt.Sprintf(" ~ %.2f%%", stdHigherCatchProb)
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Top 4",
			Value:  strings.Join(topFourDesp, "\n"),
			Inline: true,
		},
		{
			Name:   "Standard",
			Value:  "Standard Balls (Poke/Luxury/Premier): `" + standardDescription + "`",
			Inline: true,
		},
	}

	return embed, nil
}

// getBestBallsForPokemon returns
func (b *Bot) getBestBallsForPokemon(pkm *pokemon, isGmax bool) []*pokeBall {
	balls := make([]*pokeBall, 0)
	for _, pb := range b.pokemonRepo.ballsCatchRatesForPokemon(pkm) {
		if isExcludedBall(pb) {
			continue
		}
		balls = append(balls, pb)
	}

	return balls[:4]
}

func getFormFromArgs(args []string) string {

	// the form will always be the second element
	// e.g: blastoise gmax
	if len(args) < 2 || args[1] == "" {
		return ""
	}

	return getSpriteForm(args[1])
}

func (b *Bot) getPokemonColor(pType string) int {
	t, err := b.pokemonRepo.pokemonType(pType)
	if err != nil {
		return b.config.Bot.EmbedColor
	}
	return t.Color
}

func confidenceEmoji(confidence bool) string {
	if confidence {
		return "✅"
	}
	return "⛔"
}
