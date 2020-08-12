package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

const (
	galarian = "galarian"
	alolan   = "alolan"

	// Gigantamax is the gmax form
	Gigantamax = "gigantamax"
)

var (

	// ErrDenDoesNotExist is returned when a den does not exist
	ErrDenDoesNotExist = errors.New("den does not exist")

	// ErrBallDoesNotExist is returned when a ball does not exist
	ErrBallDoesNotExist = errors.New("ball does not exist")

	// ErrPokemonDoesNotExist is returned when a pokemon does not exist
	ErrPokemonDoesNotExist = errors.New("pokemon does not exist")

	// ErrTypeDoesNotExist is returned when a type does not exist
	ErrTypeDoesNotExist = errors.New("type does not exist")
)

// Repository handles all the searching and saving of all of storage related
// data in the app.
type Repository struct {
	db      *gorm.DB
	cache   *cache.Cache
	dens    map[string]*Den
	balls   map[string]*PokeBall
	pokemon map[string]*Pokemon
	types   map[string]*PokemonType
}

// NewRepository creates a new instance of the repository
// This method will load up the json files inside the /data
// folder at the project's root, and create maps for quick look up.
//
// TODO: This storage repo is a hybrid of JSON + Postgres for now till
//  everything is migrated over
func NewRepository(db *gorm.DB) (*Repository, error) {
	// load all the json files, starting by dens
	dens := make([]*Den, 0)
	if err := loadJSONInto("data/dens.json", &dens); err != nil {
		return nil, fmt.Errorf("failed to load dens.json: %+v", err)
	}

	// build the dens map for quick lookups
	densMap := make(map[string]*Den)
	for _, den := range dens {
		// if the file is structured properly, this should just work...
		// todo(hector) - add validation here to prevent panics or overwrites?
		densMap[den.Number] = den
	}

	balls := make([]*PokeBall, 0)
	if err := loadJSONInto("data/balls.json", &balls); err != nil {
		return nil, fmt.Errorf("failed to load balls.json: %+v", err)
	}

	ballsMap := make(map[string]*PokeBall)
	for _, ball := range balls {
		// all names are "<something> Ball", so just remove the " Ball" part
		lowered := strings.ToLower(ball.Name)
		ball.ID = strings.ReplaceAll(lowered, " ball", "")
		ballsMap[ball.ID] = ball
	}

	pokemons := make([]*Pokemon, 0)
	if err := loadJSONInto("data/pokemon.json", &pokemons); err != nil {
		return nil, fmt.Errorf("failed to load pokemon.json: %+v", err)
	}

	pkmMap := make(map[string]*Pokemon)
	for _, pkm := range pokemons {
		lowered := strings.ToLower(pkm.Name)
		pkmMap[lowered] = pkm
	}

	types := make([]*PokemonType, 0)
	if err := loadJSONInto("data/types.json", &types); err != nil {
		return nil, fmt.Errorf("failed to load types.json: %+v", err)
	}

	typesMap := make(map[string]*PokemonType)
	for _, pkmType := range types {
		lowered := strings.ToLower(pkmType.Name)
		typesMap[lowered] = pkmType
	}

	return &Repository{
		dens:    densMap,
		balls:   ballsMap,
		pokemon: pkmMap,
		types:   typesMap,
		db:      db,
		cache:   cache.New(5*time.Minute, 10*time.Minute),
	}, nil
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
