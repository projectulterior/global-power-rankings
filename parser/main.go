package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	GAMES_DIRECTIORY = "../data/games"
	GAMES_FILE       = "../data/games.json"
	MAPPING_FILE     = "../data/mapping_data.json"
	OUTPUT_PATH      = "../data/games_ts.json"
)

func main() {
	run()
}

func getGamePath(gameID string) string {
	return fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID)
}

func parse[T any](path string) []T {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		panic(err)
	}

	var data []T
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		panic(err)
	}

	return data
}

func write(path string, reader io.Reader) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		panic(err)
	}
}
