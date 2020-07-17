package bot

import "strings"

var (
	moonPokemon   = []string{"Nidoran", "Nidorina", "Nidoqueen", "Nidoran", "Nidorino", "Nidoking", "Cleffa", "Clefairy", "Clefable", "Igglybuff", "Jigglypuff", "Wigglytuff", "Munna", "Musharna"}
	excludedBalls = []string{"Timer", "Quick", "Master", "Level", "Lure", "Nest", "Dive", "Dream", "Heavy", "Love", "Park", "Cherish", "Sport", "Safari"}
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
