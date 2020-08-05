package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	galarian   = "galarian"
	alolan     = "alolan"
	gigantamax = "gigantamax"
)

var (
	errDenDoesNotExist     = errors.New("den does not exist")
	errBallDoesNotExist    = errors.New("ball does not exist")
	errPokemonDoesNotExist = errors.New("pokemon does not exist")
	errTypeDoesNotExist    = errors.New("type does not exist")
)

type pokemonRepo struct {
	db       *gorm.DB
	dens     map[string]*den
	balls    map[string]*pokeBall
	pokemons map[string]*pokemon
	types    map[string]*pokemonType
}

type den struct {
	Number string        `json:"den"`
	Sword  []*denPokemon `json:"sword"`
	Shield []*denPokemon `json:"shield"`
}

type denPokemon struct {
	Name       string `json:"name"`
	Ability    string `json:"ability"`
	Gigantamax bool   `json:"gigantamax"`
}

type pokeBall struct {
	ID         string
	Name       string  `json:"name"`
	Modifier   float64 `json:"modifier"`
	Conditions string  `json:"conditions"`
	Effect     string  `json:"effect"`
	Color      int     `json:"color"`
}

type pokemon struct {
	Abilities struct {
		Ability1 string `json:"ability1"`
		Ability2 string `json:"ability2"`
		AbilityH string `json:"abilityH"`
	}
	BaseStats struct {
		Atk   int `json:"atk"`
		Def   int `json:"def"`
		HP    int `json:"hp"`
		SpA   int `json:"spA"`
		SpD   int `json:"spD"`
		Spd   int `json:"spe"`
		Total int `json:"tot"`
	}
	CatchRate int `json:"catchRate"`
	Dens      struct {
		Shield []string `json:"shield"`
		Sword  []string `json:"sword"`
	}
	DexID       int      `json:"dexId"`
	EggGroup1   string   `json:"eggGroup1"`
	EggGroup2   string   `json:"eggGroup2"`
	Evolutions  []string `json:"evolutions"`
	Forms       []string `json:"forms"`
	GenderRatio string   `json:"genderRatio"`
	Generation  string   `json:"generation"`
	Height      float64  `json:"height"`
	Name        string   `json:"name"`
	Type1       string   `json:"type1"`
	Type2       string   `json:"type2"`
	Weight      float64  `json:"weight"`
	Color       int      `json:"color"`
}

type pokemonType struct {
	Name      string             `json:"name"`
	Offensive map[string]float64 `json:"offensive"`
	Defensive map[string]float64 `json:"defensive"`
	Color     int                `json:"color"`
}

// newPokemonRepo creates a new instance of the pokemonRepo
// This method will load up the json files inside the /data
// folder at the project's root, and create maps for quick look up.
func newPokemonRepo(db *gorm.DB) (*pokemonRepo, error) {
	// load all the json files, starting by dens
	dens := make([]*den, 0)
	if err := loadJSONInto("data/dens.json", &dens); err != nil {
		return nil, fmt.Errorf("failed to load dens.json: %+v:\n", err)
	}

	// build the dens map for quick lookups
	densMap := make(map[string]*den)
	for _, den := range dens {
		// if the file is structured properly, this should just work...
		// todo(hector) - add validation here to prevent panics or overwrites?
		densMap[den.Number] = den
	}

	balls := make([]*pokeBall, 0)
	if err := loadJSONInto("data/balls.json", &balls); err != nil {
		return nil, fmt.Errorf("failed to load balls.json: %+v:\n", err)
	}

	ballsMap := make(map[string]*pokeBall)
	for _, ball := range balls {
		// all names are "<something> Ball", so just remove the " Ball" part
		lowered := strings.ToLower(ball.Name)
		ball.ID = strings.ReplaceAll(lowered, " ball", "")
		ballsMap[ball.ID] = ball
	}

	pokemons := make([]*pokemon, 0)
	if err := loadJSONInto("data/pokemon.json", &pokemons); err != nil {
		return nil, fmt.Errorf("failed to load pokemon.json: %+v:\n", err)
	}

	pkmMap := make(map[string]*pokemon)
	for _, pkm := range pokemons {
		lowered := strings.ToLower(pkm.Name)
		pkmMap[lowered] = pkm
	}

	types := make([]*pokemonType, 0)
	if err := loadJSONInto("data/types.json", &types); err != nil {
		return nil, fmt.Errorf("failed to load types.json: %+v:\n", err)
	}

	typesMap := make(map[string]*pokemonType)
	for _, pkmType := range types {
		lowered := strings.ToLower(pkmType.Name)
		typesMap[lowered] = pkmType
	}

	return &pokemonRepo{
		dens:     densMap,
		balls:    ballsMap,
		pokemons: pkmMap,
		types:    typesMap,
	}, nil
}

// spriteImage returns the URL for the sprite of the pokemon. This function
// assumes the form is already validated and correct
//
// todo(hector) I hate this, I cant wait to get rid of it
func (p *pokemon) spriteImage(shiny bool, form string) string {
	fileType := "normal"
	if shiny {
		fileType = "shiny"
	}

	cleanName := strings.ReplaceAll(p.Name, "-", "")

	// the dumb mr mime line has special sprite names. I. hate. this.
	if p.DexID == 439 || p.DexID == 122 || p.DexID == 866 {
		cleanName = strings.ReplaceAll(p.Name, " ", "-")
	}

	// farfetch'd also has funny sprites _flipstable_
	if p.DexID == 83 || p.DexID == 865 {
		cleanName = strings.ReplaceAll(p.Name, "'", "")
	}

	// urshifu be special too, you thought it ended with farfetch'd? no sir!
	// everything is stupid
	if p.DexID == 892 {
		cleanName = strings.ReplaceAll(p.Name, " ", "-")
		cleanName = strings.ReplaceAll(cleanName, "-Strike-Style", "")
	}
	return fmt.Sprintf(
		"https://raw.githubusercontent.com/caquillo07/rotom-b-data/master/sprites/pokemon/%s/%s.gif",
		fileType,
		spriteFileName(strings.ToLower(cleanName), strings.ToLower(form)),
	)
}

// captureRate returns the catch rate and confidence level for the given
// poke ball, and stats combination.
func (p *pokemon) captureRate(ball *pokeBall, level, iv int, isGmax, isPromo bool) (float64, bool) {
	hpStat := hpStatFromBase(p.BaseStats.HP, iv, level)
	pCatchRate := p.CatchRate
	confidence := true

	if isGmax {
		pCatchRate = 3
	} else if isPromo {
		pCatchRate = 20
	}

	// if the modified catch rate is above 200, the confidence level drops
	catchRate := modifiedCatchRate(hpStat, 1, pCatchRate, ball.Modifier)

	// if the catch is is over 200, due to a rounding error we are not confident
	// the calculation is very accurate
	// https://bulbapedia.bulbagarden.net/wiki/Catch_rate#Probability_of_capture
	if catchRate > 200 {
		confidence = false
	}
	// if the catch rate is 255, then its a guaranteed catch. no need to keep
	// processing
	if catchRate == 255 {
		return 100, confidence
	}

	shakeProb := shakeProbability(catchRate)
	catchProb := catchProbability(shakeProb)

	return catchProb, confidence
}

// den will try to find the given den, if it does not exist it
// will return a `errDenDoesNotExist` error
func (r *pokemonRepo) den(denNumber string) (*den, error) {
	if d, ok := r.dens[denNumber]; ok {
		// return a copy of the original to not have unwanted changes to map
		// storage
		c := *d
		return &c, nil
	}
	return nil, errDenDoesNotExist
}

// ball will try to find the given ball, if it does not exist it
// will return a `errBallDoesNotExist` error
func (r *pokemonRepo) ball(ball string) (*pokeBall, error) {

	// first clean the input a bit, we will remove all white spaces,
	// then remove the "ball" part and any hyphens/underscores if present,
	// and account for any other names the ball may go for.
	ball = strings.ToLower(ball)
	for _, mod := range []string{"ball", "-", "_"} {
		ball = strings.ReplaceAll(ball, mod, "")
	}
	ball = strings.TrimSpace(ball)
	switch ball {
	case "lux":
		ball = "luxury"
	default:
		// ignore
	}

	if b, ok := r.balls[ball]; ok {
		// return a copy of the original to not have unwanted changes to map
		// storage
		c := *b
		return &c, nil
	}
	return nil, errBallDoesNotExist
}

// ballsCatchRatesForPokemon returns a list of all poke balls, sorted by
// catch effectiveness for the given pokemon.
func (r *pokemonRepo) ballsCatchRatesForPokemon(pkm *pokemon) []*pokeBall {
	// make a copy of the balls so we don't accidentally modify the global
	// state. This array is very small, so its ok to do this on every command
	// call that needs it.
	balls := make([]*pokeBall, 0)
	for _, ball := range r.balls {
		// make a copy of the ball, in case we need to modify it
		newBall := *ball
		newBall.Modifier = newBall.catchModifier(pkm)
		balls = append(balls, &newBall)
	}

	// sort them by highest modifier first
	sort.Slice(balls, func(i, j int) bool {
		return balls[i].Modifier > balls[j].Modifier
	})
	return balls
}

// pokemon will try to find the given pokemon, if it does not exist it
// will return a `errBallDoesNotExist` error
func (r *pokemonRepo) pokemon(name string) (*pokemon, error) {
	if p, ok := r.pokemons[strings.ToLower(name)]; ok {
		// return a copy of the original to not have unwanted changes to map
		// storage
		c := *p
		return &c, nil
	}
	return nil, errPokemonDoesNotExist
}

// pokemon will try to find the given pokemon, if it does not exist it
// will return a `errBallDoesNotExist` error
func (r *pokemonRepo) pokemonType(name string) (*pokemonType, error) {
	if t, ok := r.types[strings.ToLower(name)]; ok {
		// return a copy of the original to not have unwanted changes to map
		// storage
		c := *t
		return &c, nil
	}
	return nil, errTypeDoesNotExist
}

// catchModifier returns the actual modifier for the given pokemon
func (b *pokeBall) catchModifier(pkm *pokemon) float64 {
	mod := b.Modifier
	switch b.ID {
	case "moon":
		mod = 1
		if isMoonPokemon(pkm) {
			mod = 3.5
		}
	case "fast":
		mod = 1
		if pkm.BaseStats.SpD >= 100 {
			mod = 4.0
		}
	case "net":
		mod = 1
		if pkm.Type1 == "Water" || pkm.Type1 == "Bug" || pkm.Type2 == "Water" || pkm.Type2 == "Bug" {
			mod = 3.5
		}
	case "love":
		mod = 1
		if pkm.GenderRatio != "100% ⚲" && pkm.GenderRatio != "100% ♀" && pkm.GenderRatio != "100% ♂" {
			mod = 8
		}
	case "heavy":
		mod = 1
		if pkm.Weight >= 300 {
			mod = 30
		} else if pkm.Weight >= 200 {
			mod = 20
		}
	case "beast":
		mod = 0.1
	default:
		// do nothing
	}
	return mod
}

func loadJSONInto(fileLocation string, i interface{}) error {
	denFile, err := os.Open(fileLocation)
	if err != nil {
		return err
	}

	densBuf, err := ioutil.ReadAll(denFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(densBuf, &i)
}

// I think this was a waste of time, but will leave because I spent so long
// on it :cries:
// nolint:goconst
func getSpriteForm(form string) string {
	// lets make it lower case to make our life's easier. If we get the full
	// form name, just return the capitalized version of it. If we get a
	// shorthand name, return a static string
	f := strings.ToLower(form)
	switch f {
	case "galarian", "alolan", "gigantamax", "mega", "pirouette", "busted",
		"midnight", "dusk", "hangry", "pau", "sensu", "pompom", "cosplay",
		"primal", "blade-form", "sunny", "rainy", "snowy", "sunshine",
		"crowned", "starter", "black", "white", "f", "dusk-mane", "dawn-wings",
		"megay", "megax", "gorging", "gulping", "dusk-form", "zen-mode",
		"galarian-zen", "ultra", "origin-form":
		return f
	case "galar":
		return "galarian"
	case "alola":
		return "Alolan"
	case "gmax", "g-max":
		return "gigantamax"
	case "mega-y", "mega-x":
		return strings.ReplaceAll(f, "-", "")
	case "pom-pom":
		return "pompom"
	case "female":
		return "f"
	case "blade":
		return "blade-form"
	case "zen":
		return "zen-mode"
	case "galar-zen", "galar-zen-mode", "galarian-zen-mode":
		return "galarian-zen"
	case "origin", "origin-forme":
		return "origin-form"
	// this cases don't have special naming for sprites, so just fall
	// through and return
	case "aria", "m", "male", "disguised":
		fallthrough
	default:
		// nothing
	}
	return ""
}

// todo(hector) remove this non-sense as soon as images re-formaterd
func spriteFileName(pkm, form string) string {
	switch form {
	case alolan, "crowned", "dusk", "midnight", galarian, "mega", "megay",
		"megax", "primal", "ultra":
		return fmt.Sprintf("%s-%s", form, pkm)
	case "":
		return pkm
	case "dusk-form":
		return fmt.Sprintf("dusk-%s", pkm)
	case "galarian-zen":
		return fmt.Sprintf("galarian-%s-zen", pkm)
	default:
		// ignore
	}
	return fmt.Sprintf("%s-%s", pkm, form)
}
