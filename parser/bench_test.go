package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	simd "github.com/minio/simdjson-go"
)

const (
	GAMES_COUNT = 100
)

func BenchmarkStandardJSON(b *testing.B) {
	benchmarkTest(jsonControl)
}

func BenchmarkSIMDJSON(b *testing.B) {
	simdBenchmarkTest()
}

func simdBenchmarkTest() {
	games := parse[map[string]any](GAMES_FILE)

	var wg sync.WaitGroup
	sema := make(chan struct{}, 30)

	var count atomic.Int32

	start := time.Now()
	for i, g := range games {
		if i >= GAMES_COUNT {
			break
		}

		wg.Add(1)
		sema <- struct{}{}
		go func(game map[string]any) {
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

			simdAnal(*parsed)
		}(g)
	}
	wg.Wait()

	fmt.Printf("errors: %d -- %s\n", count.Load(), time.Since(start))
	fmt.Println("------------")
}

func benchmarkTest(unmarshalFn func(any, bytes.Buffer) error) {
	games := strategicParse[map[string]any](GAMES_FILE, unmarshalFn)

	var wg sync.WaitGroup
	sema := make(chan struct{}, 30)

	var count atomic.Int32

	start := time.Now()
	for i, g := range games {
		if i >= GAMES_COUNT {
			break
		}

		wg.Add(1)
		sema <- struct{}{}
		go func(game map[string]any) {
			defer wg.Done()
			defer func() { <-sema }()

			gameID := game["id"].(string)

			data := strategicGetGame(gameID, unmarshalFn)
			if len(data) == 0 {
				fmt.Println("no file found")
				return
			}

			analyze(data)
		}(g)
	}
	wg.Wait()

	fmt.Printf("errors: %d -- %s\n", count.Load(), time.Since(start))
	fmt.Println("------------")
}

/* JSON decode algorithms */
func jsonControl(data any, buffer bytes.Buffer) error {
	err := json.Unmarshal(buffer.Bytes(), data)
	if err != nil {
		return err
	}

	return nil
}

func jsonJSONIter(data any, buffer bytes.Buffer) error {
	err := jsoniter.Unmarshal(buffer.Bytes(), data)
	if err != nil {
		return err
	}

	return nil
}

func jsonSIMD(data any, buffer bytes.Buffer) error {
	parsedJSON, err := simd.Parse(buffer.Bytes(), nil)
	if err != nil {
		return err
	}

	data = parsedJSON
	return nil
}
