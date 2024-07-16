package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"pokedex.ayonchakroborty.net/internals/pokeapi"
)

type cliCommand struct{
	name		string
	description	string
	callBack	func() error
}

func startRepl(){
	scanner := bufio.NewScanner(os.Stdin)
	commands := newCommandsMap()
	
	pokeMap := pokeapi.PokedexMap{}
	if entry, ok := commands["map"]; ok{
		entry.callBack = pokeMap.CommandMapf
		commands["map"] = entry
	}

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())

		if exitRepl(input){
			break
		}

		err := commands[input[0]].callBack()
		if err != nil{
			panic(err)
		}

		fmt.Println()
	}
}

func cleanInput(input string) []string{
	input = strings.ToLower(input)
	return strings.Fields(input)
}

func exitRepl(input []string) bool{
	return len(input) == 1 && strings.Compare(input[0], "exit") == 0
}

func newCommandsMap() map[string]cliCommand{
	return map[string]cliCommand{
		"help" : {
			name:	"help",
			description: "Dispays a help message",
			callBack: commandHelp,
		},
		"exit" : {
			name:	"exit",
			description: "Exits the Pokedex",
			callBack: commandExit,
		},
		"map":{
			name: "map",
			description: "Get the next page of locations",
		},
		"mapb":{
			name: "mapb",
			description: "Get the previous page of locations",
		},
	}
}