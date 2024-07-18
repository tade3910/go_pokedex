package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	pokeApi "github.com/tade3910/go_pokedex/internal/pokeApi"
)

func printHelp() {
	defer fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	fmt.Println("help: Displays a help message")
	fmt.Println("catch some_pokemon: Attempts to catch some_pokemon")
	fmt.Println("explore some_location: Explores some_location and displays pokemen there")
	fmt.Println("map: Displays the names of the next 20 location areas in the Pokemon world")
	fmt.Println("mapb: Displays the previous 20 locations")
	fmt.Println("pokedex: Prints a list of all the names of the Pokemon the user has caught")
	fmt.Println("exit: Exit the Pokedex")
}

func handleMap(prev bool, pokeConfig *pokeApi.Config) {
	res, err := pokeConfig.GetMap(prev)
	if err != nil {
		log.Fatal(err)
	}
	if res.Next != nil {
		pokeConfig.Next = *res.Next
	} else {
		pokeConfig.Next = ""
	}
	if res.Previous != nil {
		pokeConfig.Previous = *res.Previous
	} else {
		pokeConfig.Previous = ""
	}

	for _, result := range res.Results {
		fmt.Println(result.Name)
	}
}

func handleExplore(pokeConfig *pokeApi.Config, key string) {
	fmt.Printf("Exploring %s...\n", key)
	area, err := pokeConfig.GetLocation(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon_encounter := range area.PokemonEncounters {
		name := pokemon_encounter.Pokemon.Name
		fmt.Printf("- %s\n", name)
	}
}

func handleCatch(pokeConfig *pokeApi.Config, key string) {
	fmt.Printf("Throwing a Pokeball at %s...\n", key)
	pokemon, err := pokeConfig.GetPokemon(key)
	if err != nil {
		log.Fatal(err)
	}
	chance := 180 * rand.Float64()
	fmt.Printf("chance to catch is %f, base experience of %s is %d\n", chance, key, pokemon.BaseExperience)
	if chance > float64(pokemon.BaseExperience) {
		pokeConfig.PokeDex[key] = pokemon
		fmt.Printf("%s was caught!\n", key)
	} else {
		fmt.Printf("%s escaped\n", key)
	}
}

func handleInspect(pokeConfig *pokeApi.Config, key string) {
	pokemon, caught := pokeConfig.PokeDex[key]
	if !caught {
		fmt.Printf("You have not yet caught %s\n", key)
		return
	}
	pokemon.PrintDetails()
}

func handleError() {
	fmt.Println("Unknown command")
}

func handlePokeDex(pokeConfig *pokeApi.Config) {
	fmt.Println("Your Pokedex:")
	for pokemon := range pokeConfig.PokeDex {
		fmt.Printf(" - %s\n", pokemon)
	}
}

func handleCommand(command string, pokeConfig *pokeApi.Config) {
	switch command {
	case "help":
		printHelp()
	case "map":
		handleMap(false, pokeConfig)
	case "mapb":
		handleMap(true, pokeConfig)
	case "pokedex":
		handlePokeDex(pokeConfig)
	default:
		handleError()
	}
}

func handle2Command(command string, identifier string, pokeConfig *pokeApi.Config) {
	switch command {
	case "explore":
		handleExplore(pokeConfig, identifier)
	case "catch":
		handleCatch(pokeConfig, identifier)
	case "inspect":
		handleInspect(pokeConfig, identifier)
	default:
		handleError()
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	pokeConfig := pokeApi.GetNewConfig()
	for {
		fmt.Print("Pokedex > ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		words := strings.Fields(text)
		lenth := len(words)
		switch lenth {
		case 2:
			handle2Command(words[0], words[1], &pokeConfig)
		case 1:
			if text == "exit" {
				return
			}
			handleCommand(text, &pokeConfig)
		default:
			handleError()
		}
	}
}
