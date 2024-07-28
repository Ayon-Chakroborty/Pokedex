package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type Entry struct {
	data   []byte
	caught bool
}
type Pokedex struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	PokedexEntries map[string]Entry
}

func (p *Pokedex) getData(pokemonName string) ([]byte, error) {
	if body, ok := p.PokedexEntries[pokemonName]; ok && body.data != nil {
		return body.data, nil
	}

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemonName)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Compare(string(body), "Not Found") == 0 {
		fmt.Printf("\nError: %s is not a real pokemon\n\n", pokemonName)
		return nil, nil
	}

	p.PokedexEntries[pokemonName] = Entry{
		data:   body,
		caught: false,
	}
	return body, nil
}

func (p *Pokedex) Catch(pokemonName string) error {
	if val, ok := p.PokedexEntries[pokemonName]; ok && val.caught {
		fmt.Printf("\n%s was already caught\n\n", pokemonName)
		return nil
	}

	body, err := p.getData(pokemonName)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return err
	}

	chance := rand.Intn(p.BaseExperience)
	fmt.Printf("\nThrowing a Pokeball at %s...\n", p.Name)
	if chance > 40 {
		fmt.Printf("%s escaped!\n\n", p.Name)
	} else {
		fmt.Printf("%s was caught!\n\n", p.Name)
		p.PokedexEntries[pokemonName] = Entry{body, true}
	}

	return nil
}

func (p *Pokedex) Inspect(pokemonName string) error {
	body, err := p.getData(pokemonName)
	if err != nil {
		return nil
	}
	if body == nil {
		return nil
	}

	if entry, ok := p.PokedexEntries[pokemonName]; !ok || (ok && !entry.caught){
		fmt.Println("\nyou have not caught that pokemon\n")
		return nil
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf("	- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range p.Types {
		fmt.Printf(" - %s\n", pokemonType.Type.Name)
	}
	fmt.Println()
	return nil
}
