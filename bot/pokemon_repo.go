package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	errDenDoesNotExist     = errors.New("den does not exist")
	errBallDoesNotExist    = errors.New("ball does not exist")
	errPokemonDoesNotExist = errors.New("pokemon does not exist")
)

type pokemonRepo struct {
	dens     map[string]*den
	balls    map[string]*pokeBall
	pokemons map[string]*pokemon
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
}

// newPokemonRepo creates a new instance of the pokemonRepo
// This method will load up the json files inside the /data
// folder at the project's root, and create maps for quick look up.
func newPokemonRepo() (*pokemonRepo, error) {
	// load all the json files, starting by dens
	dens := make([]*den, 0)
	if err := loadJSONInto("data/dens.json", &dens); err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	pkmMap := make(map[string]*pokemon)
	for _, pkm := range pokemons {
		lowered := strings.ToLower(pkm.Name)
		pkmMap[lowered] = pkm
	}
	return &pokemonRepo{
		dens:     densMap,
		balls:    ballsMap,
		pokemons: pkmMap,
	}, nil
}

func (p *pokemon) isGigantamax() bool {
	for _, f := range p.Forms {
		if f == "Gigantamax" {
			return true
		}
	}
	return false
}

// spriteImage returns the URL for the sprite of the pokemon
//
// todo(hector) I hate this, I cant wait to get rid of it
func (p *pokemon) spriteImage(shiny bool, form string) string {
	fileName := strings.ToLower(p.Name)
	fileType := "normal"
	if shiny {
		fileType = "shiny"
	}

	if form != "" {
		fileName = fmt.Sprintf("%s-%s", fileName, strings.ToLower(form))
	}

	return fmt.Sprintf(
		"https://raphgg.github.io/den-bot/data/sprites/pokemon/%s/%s.gif",
		fileType,
		fileName,
	)
}

// den will try to find the given den, if it does not exist it
// will return a `errDenDoesNotExist` error
func (r *pokemonRepo) den(denNumber string) (*den, error) {
	if d, ok := r.dens[denNumber]; ok {
		return d, nil
	}
	return nil, errDenDoesNotExist
}

// ball will try to find the given ball, if it does not exist it
// will return a `errBallDoesNotExist` error
func (r *pokemonRepo) ball(ball string) (*pokeBall, error) {
	if b, ok := r.balls[strings.ToLower(ball)]; ok {
		return b, nil
	}
	return nil, errBallDoesNotExist
}

// pokemon will try to find the given pokemon, if it does not exist it
// will return a `errBallDoesNotExist` error
func (r *pokemonRepo) pokemon(name string) (*pokemon, error) {
	if p, ok := r.pokemons[strings.ToLower(name)]; ok {
		return p, nil
	}
	return nil, errPokemonDoesNotExist
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
