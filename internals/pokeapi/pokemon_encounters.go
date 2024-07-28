package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	cache "pokedex.ayonchakroborty.net/internals/pokecache"
)

type Pokemons struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
	GameCache *cache.PokeCache
}

var baseUrl = "https://pokeapi.co/api/v2/location-area/"

func (p *Pokemons) getData(area string)([]byte, error){
	if val, exists := p.GameCache.GetPokemon(area); exists{
		return val, nil
	}

	resp, err := http.Get(baseUrl+area)
	defer resp.Body.Close()
	if err != nil{
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	p.GameCache.AddPokemon(area, body)
	return body, nil
}

func (p *Pokemons) Explore(area string) error {
	body, err := p.getData(area)
	if err != nil{
		return err
	}
	if strings.Compare(string(body), "Not Found") == 0{
		fmt.Printf("\nError: %s is not a real area\n\n", area)
		return nil
	}

	err = json.Unmarshal(body, p)
	if err != nil{
		return err
	}

	fmt.Println()
	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon:")
	for _, v := range p.PokemonEncounters{
		fmt.Printf("- %s\n", v.Pokemon.Name)
	}

	fmt.Println()
	return nil
}
