package bot

import "math"

// All calculations here were made following the formulas provided by bulbapedia
// base stats: https://bulbapedia.bulbagarden.net/wiki/Statisticz
//
// catch rates: https://bulbapedia.bulbagarden.net/wiki/Catch_rate#Modified_catch_rate_3

func hpStatFromBase(baseHP, iv, level int) int {
	x := ((2 * float64(baseHP)) + float64(iv)) * float64(level)
	return int(math.Floor(x/100) + float64(level) + 10)
}

func modifiedCatchRate(maxHP, currentHP, catchRate int, ballModifier float64) float64 {
	x := ((3 * maxHP) - (2 * currentHP)) * catchRate
	return (float64(x) * ballModifier) / float64(3*maxHP)
}

func shakeProbability(modifiedCatchRate float64) float64 {
	return math.Floor(65536 / math.Pow(255/modifiedCatchRate, 0.1875))
}

func catchProbability(shakeProb float64) float64 {
	if shakeProb >= 65536 {
		return 100
	}
	return math.Pow(shakeProb/65535, 4) * 100
}
