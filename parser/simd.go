package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	simd "github.com/minio/simdjson-go"
)

const WORKERS = 30

func run() {
	games := parse[map[string]any](GAMES_FILE)

	var wg sync.WaitGroup
	sema := make(chan struct{}, WORKERS)

	var count atomic.Int32

	// data := make(chan *Game, WORKERS)
	// go func() {
	// 	toRet := []map[string]any{}
	// 	for game := range data {
	// 		currentGame := games[i]
	// 		game["start_time"] = gameTimes.Start
	// 		game["end_time"] = gameTimes.End
	// 	}

	// 	b, err := json.Marshal(games)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	write(OUTPUT_PATH, bytes.NewBuffer(b))
	// }()

	start := time.Now()
	for i, g := range games {
		wg.Add(1)
		sema <- struct{}{}
		go func(index int, game map[string]any) {
			defer wg.Done()
			defer func() { <-sema }()

			gameID := game["id"].(string)

			parsed, err := simdParse(fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID))
			if parsed == nil {
				fmt.Printf("parsed is null -- %s\n", err)
				return
			}

			iter := parsed.Iter()
			if iter.PeekNext() == simd.TypeNone {
				fmt.Println("no file found")
				return
			}

			gameTimes, err := simdAnal(*parsed)
			if err != nil {
				fmt.Printf("analysis error -- %s\n", err)
			}

			currentGame := games[index]
			currentGame["start_time"] = gameTimes.Start
			currentGame["end_time"] = gameTimes.End
		}(i, g)

		if i%100 == 0 {
			fmt.Printf("checkpoint -- %d -- %s\n", i, time.Since(start))
		}
	}
	wg.Wait()
	// close(data)

	b, err := json.Marshal(games)
	if err != nil {
		panic(err)
	}

	write(OUTPUT_PATH, bytes.NewBuffer(b))
	fmt.Printf("errors: %d -- %s\n", count.Load(), time.Since(start))
	fmt.Println("------------")
}

func strategicGetGame(gameID string, unmarshalFn func(any, bytes.Buffer) error) Events {
	return strategicParse[Event](fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID), unmarshalFn)
}

// Parse functions
func strategicParse[T any](path string, unmarshalFn func(any, bytes.Buffer) error) []T {
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
	err = unmarshalFn(&data, buf)
	if err != nil {
		fmt.Printf("error unmarshaling -- %s", err)
	}

	return data
}

func simdParse(path string) (*simd.ParsedJson, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	parsed, err := simd.Parse(buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
