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
	mapping := parse[map[string]any](MAPPING_FILE)

	var wg sync.WaitGroup
	sema := make(chan struct{}, WORKERS)

	var count atomic.Int32

	start := time.Now()
	for i, g := range games[:100] {
		wg.Add(1)
		sema <- struct{}{}
		go func(index int, game map[string]any, mapping []map[string]any) {
			defer wg.Done()
			defer func() { <-sema }()

			gameID := game["id"].(string)

			parsed, err := simdParse(getGamePath(gameID))
			if parsed == nil {
				fmt.Printf("parsed is null -- %s\n", err)
				return
			}

			iter := parsed.Iter()
			if iter.PeekNext() == simd.TypeNone {
				fmt.Println("no file found")
				return
			}

			gameTimes, err := analyze(*parsed)
			if err != nil {
				fmt.Printf("analysis error -- %s\n", err)
			}

			currentGame := games[index]
			currentGame["start_time"] = gameTimes.Start
			currentGame["end_time"] = gameTimes.End
		}(i, g, mapping)

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
