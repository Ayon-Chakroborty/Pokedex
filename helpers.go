package main

import (
	"fmt"
	"strings"
)

func commandHelp(input string) error{
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usages:\n")
	for _, v := range newCommandsMap(){
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func commandExit(input string) error{
	return nil
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

	if entry, ok := c["catch"]; ok {
		entry.callBack = app.pokedex.Catch
		c["catch"] = entry
	}

	if entry, ok := c["inspect"]; ok {
		entry.callBack = app.pokedex.Inspect
		c["inspect"] = entry
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