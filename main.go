package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	pokemap "github.com/tade3910/go_pokedex/internal/pokeMap"
)

func printHelp() {
	defer fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	fmt.Println("map: Displays the names of the next 20 location areas in the Pokemon world")
	fmt.Println("mapb: Displays the previous 20 locations")
}

func handleMap(prev bool, pokeConfig *pokemap.Config) {
	res, err := pokeConfig.GetMap(prev)
	if err != nil {
		log.Fatal(err)
	}
	if res.Next != nil {
		pokeConfig.Next = *res.Next
	}
	if res.Previous != nil {
		pokeConfig.Previous = *res.Previous
	}

	for _, result := range res.Results {
		fmt.Println(result.Name)
	}
}

func handleCommand(command string, pokeConfig *pokemap.Config) {
	switch command {
	case "help":
		printHelp()
	case "map":
		handleMap(false, pokeConfig)
	case "mapb":
		handleMap(true, pokeConfig)
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	pokeConfig := pokemap.GetNewConfig()
	for {
		fmt.Print("Pokedex > ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text == "exit" {
			break
		}
		handleCommand(text, &pokeConfig)
	}
}
