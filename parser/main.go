package main

import "os"

const (
	GAMES_DIRECTIORY = "data/games/"
	GAMES_FILE       = "data/games.json"
)

func main() {
	games, err := os.Open(GAMES_FILE)
	if err != nil {
		panic(err)
	}
}
