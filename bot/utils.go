package bot

import "strings"

var (
	moonPokemon = []string{
		"Nidoran", "Nidorina", "Nidoqueen", "Nidoran", "Nidorino", "Nidoking",
		"Cleffa", "Clefairy", "Clefable", "Igglybuff", "Jigglypuff",
		"Wigglytuff", "Munna", "Musharna",
	}
	excludedBalls = []string{"Timer", "Quick", "Master", "Level", "Lure",
		"Nest", "Dive", "Dream", "Heavy", "Love", "Park", "Cherish", "Sport",
		"Safari",
	}

	ballNames = []string{"poke", "great", "ultra", "premier", "luxury", "lux",
		"beast", "cherish", "dive", "dream", "dusk", "fast", "friend", "heal",
		"heavy", "level", "love", "lure", "master", "moon", "nest", "net",
		"park", "quick", "repeat", "safari", "sport", "timer",
	}

	natures = map[string]string{
		"Hardy":   "No changes", // nolint:goconst
		"Lonely":  "+Atk -Def",
		"Brave":   "+Atk -Spe",
		"Adamant": "+Atk -SpA",
		"Naughty": "+Atk -SpD",
		"Bold":    "+Def -Atk",
		"Docile":  "No changes",
		"Relaxed": "+Def -Spe",
		"Impish":  "+Def -SpA",
		"Lax":     "+Def -SpD",
		"Timid":   "+Spe -Atk",
		"Hasty":   "+Spe -Def",
		"Serious": "No changes",
		"Jolly":   "+Spe -SpA",
		"Naive":   "+Spe -SpD",
		"Modest":  "+SpA -Atk",
		"Mild":    "+SpA -Def",
		"Quiet":   "+SpA -Spe",
		"Bashful": "No changes",
		"Rash":    "+SpA -SpD",
		"Calm":    "+SpD -Atk",
		"Gentle":  "+SpD -Def",
		"Sassy":   "+SpD -Spe",
		"Careful": "+SpD -SpA",
		"Quirky":  "No changes",
	}
)

func isMoonPokemon(pkm *pokemon) bool {
	for _, name := range moonPokemon {
		if strings.EqualFold(pkm.Name, name) {
			return true
		}
	}
	return false
}

func isExcludedBall(ball *pokeBall) bool {
	for _, b := range excludedBalls {
		if strings.EqualFold(ball.ID, b) {
			return true
		}
	}
	return false
}
