package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cache "pokedex.ayonchakroborty.net/internals/pokecache"
)

type PokedexMap struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Names    []struct {
		Name string `json:"name"`
	} `json:"results"`
	GameCache *cache.PokeCache
}

func (p *PokedexMap) initMap(url *string) error {
	if (p.Next != nil || p.Previous != nil) && url == nil {
		return nil
	}
	var (
		body   []byte
		err    error
		exists bool
	)

	if url != nil {
		body, exists = p.GameCache.GetLocation(*url)
	}

	if !exists {
		body, err = p.getData(url)
	}

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *PokedexMap) getData(url *string) ([]byte, error) {
	var (
		resp *http.Response
		err  error
	)

	if p.Next == nil && p.Previous == nil {
		urlStr := "https://pokeapi.co/api/v2/location-area/"
		url = &urlStr
	}

	resp, err = http.Get(*url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p.GameCache.AddLocation(*url, body)
	return body, nil
}

func (p *PokedexMap) CommandMapf(input string) error {
	err := p.initMap(p.Next)
	if err != nil {
		return nil
	}

	fmt.Println()
	for _, name := range p.Names {
		fmt.Printf("%s\n", name.Name)
	}

	return nil
}

func (p *PokedexMap) CommandMapb(input string) error {
	err := p.initMap(p.Previous)
	if err != nil {
		return nil
	}

	fmt.Println()
	for _, name := range p.Names {
		fmt.Printf("%s\n", name.Name)
	}

	return nil
}
