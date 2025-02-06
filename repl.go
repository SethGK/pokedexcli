package main

import (
	"strings"
	"bufio"
	"os"
	"fmt"
)

// REPL 
func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		fmt.Printf("Your command was: %s\n", words[0])
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	return words
}


	