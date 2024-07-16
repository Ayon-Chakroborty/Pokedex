package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokedexMap struct{
	Next		*string		`json:"next"`
	Previous	*string		`json:"previous"`
	Names		[]struct{
		Name	string		`json:"name"`
		url		string		`json:"url"`
	} 						`json:"results"`

}

func (p *PokedexMap) initMap(direction *string) error{
	var(
		resp 	*http.Response
		err 	error
	)
	
	if p.Next == nil && p.Previous == nil{
		resp, err = http.Get("https://pokeapi.co/api/v2/location-area/")
	} else {
		resp, err = http.Get(*direction)
	}

	if err != nil{
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return err
	}

	err = json.Unmarshal(body, p)
	if err != nil{
		return err
	}

	return nil
}

func (p *PokedexMap) CommandMapf() error{
	err := p.initMap(p.Next)
	if err != nil{
		return nil
	}

	for _, name := range p.Names{
		fmt.Printf("%s\n", name.Name)
	}

	return nil
}

func CommandMapb() error{
	resp, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil{
		return err
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Println(body)
	if err != nil{
		return err
	}

	return nil
}