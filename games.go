package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	GAMES_PATH   = "data/games.json"
	MAPPING_PATH = "data/mapping_data.json"

	S3 = "https://power-rankings-dataset-gprhack.s3.us-west-2.amazonaws.com"

	OUTPUT_PATH = "data/games"

	WORKERS = 10
)

func main() {
	games := parse(GAMES_PATH)
	mapping := parse(MAPPING_PATH)

	sema := make(chan struct{}, WORKERS)
	var wg sync.WaitGroup
	var count int
	start := time.Now()

	skipped := []string{}
	var mutex sync.Mutex

	for _, game := range games {
		wg.Add(1)
		sema <- struct{}{}
		go func(game map[string]any) {
			defer wg.Done()
			defer func() { <-sema }()

			gid, ok := game["id"].(string)
			if !ok {
				panic("error in parsing game id")
			}

			pid, err := platformID(mapping, gid)
			if err != nil {
				fmt.Printf("skipping %s: %s\n", gid, err)
				mutex.Lock()
				defer mutex.Unlock()

				skipped = append(skipped, gid)
				return
			}

			g := getGame(context.Background(), pid)
			save(fmt.Sprintf("%s/%s.json", OUTPUT_PATH, gid), g)
		}(game)

		count++

		if count%10 == 0 {
			fmt.Printf("count: %d -- %s\n", count, time.Now().Sub(start).String())
			start = time.Now()
			break
		}
	}

	wg.Wait()

	b, err := json.Marshal(skipped)
	if err != nil {
		panic(err)
	}

	save(OUTPUT_PATH+"/skipped.json", bytes.NewBuffer(b))
}

func parse(path string) []map[string]any {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		panic(err)
	}

	var data []map[string]any
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		panic(err)
	}

	return data
}

func platformID(mapping []map[string]any, gameID string) (string, error) {
	for _, game := range mapping {
		if game["esportsGameId"] == gameID {
			id, ok := game["platformGameId"].(string)
			if !ok {
				return "", fmt.Errorf("error in parsing plaform id")
			}
			return id, nil
		}
	}
	return "", fmt.Errorf("platform id not found: " + gameID)
}

func getGame(ctx context.Context, plaformID string) io.Reader {
	url := fmt.Sprintf("%s/games/%s.json.gz", S3, plaformID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(&buf, reader)
	if err != nil {
		panic(err)
	}

	return &buf
}

func save(path string, reader io.Reader) {
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
