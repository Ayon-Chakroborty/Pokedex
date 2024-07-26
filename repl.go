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
		}

		err := command.callBack("")
		if err != nil {
			panic(err)
		}

		fmt.Println()
	}
}

func addCommands(c map[string]cliCommand, app *application) {
	if entry, ok := c["map"]; ok {
		entry.callBack = app.pokeMap.CommandMapf
		c["map"] = entry
	}

	if entry, ok := c["mapb"]; ok {
		entry.callBack = app.pokeMap.CommandMapb
		c["mapb"] = entry
	}

	if entry, ok := c["explore"]; ok {
		entry.callBack = app.pokemons.Explore
		c["explore"] = entry
	}

}

func printErrorMessage(err string, helper string) {
	fmt.Printf("\n%s\n", err)
	fmt.Printf("\n%s\n\n", helper)
}

func cleanInput(input string) []string {
	input = strings.ToLower(input)
	return strings.Fields(input)
}

func exitRepl(input []string) bool {
	return len(input) == 1 && strings.Compare(input[0], "exit") == 0
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
	}
}
