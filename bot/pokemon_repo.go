package bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var (
	errDenDoesNotExist  = errors.New("den does not exist")
	errBallDoesNotExist = errors.New("ball does not exist")
)

type pokemonRepo struct {
	dens  map[string]*den
	balls map[string]*pokeBall
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
	return &pokemonRepo{
		dens:  densMap,
		balls: ballsMap,
	}, nil
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
	if b, ok := r.balls[ball]; ok {
		return b, nil
	}
	return nil, errBallDoesNotExist
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
