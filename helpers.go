package main

import (
	"fmt"
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