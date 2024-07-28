package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"pokedex.ayonchakroborty.net/internals/pokeapi"
	cache "pokedex.ayonchakroborty.net/internals/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callBack    func(input string) error
}

type application struct {
	pokeMap  *pokeapi.PokedexMap
	pokemons *pokeapi.Pokemons
	pokedex  *pokeapi.Pokedex
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := newCommandsMap()
	gameCache := cache.NewCache(time.Minute * 5)
	var errorMessage string

	app := application{
		pokeMap: &pokeapi.PokedexMap{
			GameCache: gameCache,
		},
		pokemons: &pokeapi.Pokemons{
			GameCache: gameCache,
		},
		pokedex: &pokeapi.Pokedex{
			PokedexEntries: make(map[string]pokeapi.Entry),
		},
	}

	addCommands(commands, &app)

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) > 2 {
			errorMessage = "Error: too many arguments"
			printErrorMessage(errorMessage,
				"Use 'help' command to get a list of available commands")
		}

		if exitRepl(input) {
			break
		}

		command, exist := commands[input[0]]
		if !exist {
			errorMessage = "Error: No command found with name" + input[0]
			printErrorMessage(errorMessage,
				"Use 'help' command to get a list of available commands")
			continue
		} else if strings.Compare(input[0], "explore") == 0 {
			if len(input) < 2 {
				errorMessage = "Error: too few arguments with 'explore'"
				printErrorMessage(errorMessage,
					"'explore' takes 1 arugment that represents location.\n Example: 'explore pastoria-city-area'")
				continue
			}

			err := command.callBack(input[1])
			if err != nil {
				panic(err)
			}
			continue
		} else if strings.Compare(input[0], "catch") == 0 {
			if len(input) < 2 {
				errorMessage = "Error: too few arguments with 'catch'"
				printErrorMessage(errorMessage,
					"'catch' takes 1 arugment that represents pokemon.\n Example: 'catch clefairy'")
				continue
			}

			err := command.callBack(input[1])
			if err != nil {
				panic(err)
			}
			continue
		} else if strings.Compare(input[0], "inspect") == 0 {
			if len(input) < 2 {
				errorMessage = "Error: too few arguments with 'inspect'"
				printErrorMessage(errorMessage,
					"'inspect' takes 1 arugment that represents pokemon.\n Example: 'inspect clefairy'")
				continue
			}

			err := command.callBack(input[1])
			if err != nil {
				panic(err)
			}
			continue
		}

		err := command.callBack("")
		if err != nil {
			panic(err)
		}

		fmt.Println()
	}
}

func newCommandsMap() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Dispays a help message",
			callBack:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callBack:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
		},
		"explore": {
			name:        "explore",
			description: "Explores a location and lists out the pokemon in that area",
		},

		"catch": {
			name:        "catch",
			description: "Tries to catch a specific pokemon",
		},

		"inspect": {
			name:        "inspect",
			description: "Shows details of a specific pokemon",
		},
	}
}
