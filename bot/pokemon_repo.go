package bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	errDenDoesNotExist = errors.New("den does not exist")
)

type pokemonRepo struct {
	dens map[string]*den
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

// newPokemonRepo creates a new instance of the pokemonRepo
// This method will load up the json files inside the /data
// folder at the project's root, and create maps for quick look up.
func newPokemonRepo() (*pokemonRepo, error) {
	// load the dens.json file
	denFile, err := os.Open("data/dens.json")
	if err != nil {
		return nil, err
	}

	densBuf, err := ioutil.ReadAll(denFile)
	if err != nil {
		return nil, err
	}

	dens := make([]*den, 0)
	if err := json.Unmarshal(densBuf, &dens); err != nil {
		return nil, err
	}

	// build the dens map for quick lookups
	densMap := make(map[string]*den)
	for _, den := range dens {
		// if the file is structured properly, this should just work...
		// todo(hector) - add validation here to prevent panics or overwrites?
		densMap[den.Number] = den
	}
	return &pokemonRepo{
		dens: densMap,
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
