package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/valyala/fastjson"
)

const (
	GAMES_DIRECTIORY = "../data/games"
	GAMES_FILE       = "../data/games.json"
)

func main() {
	games := parse[map[string]any](GAMES_FILE)
	fmt.Printf("length -- %d", len(games))

	var wg sync.WaitGroup
	sema := make(chan struct{}, 30)

	start := time.Now()

	var count atomic.Int32

	for i, g := range games {
		wg.Add(1)
		sema <- struct{}{}
		go func(game *fastjson.Value) {
			defer wg.Done()
			defer func() { <-sema }()

			gameID := string(game.GetStringBytes("id"))
			// defer func() {
			// 	if r := recover(); r != nil {
			// 		count.Add(1)
			// 		fmt.Printf("errors: %d -- %s\n", count.Load(), gameID)
			// 	}
			// }()

			// os.Remove(getGamePath(gameID))
			// return

			start := time.Now()

			data := getGame(gameID)
			if len(data) == 0 {
				// fmt.Println("no file found")
				return
			}

			_, err := analyze(data)
			if err != nil {
				panic(fmt.Sprintf("%s -- %s\n", gameID, err.Error()))
			}

			fmt.Printf("game analyzed -- %s\n", time.Since(start))
		}(g)

		if i%100 == 0 {
			fmt.Printf("checkpoint -- %d -- %s\n", i, time.Since(start))
		}
	}

	wg.Wait()
	fmt.Printf("errors: %d -- %s\n", count.Load(), time.Since(start))
}

func parse[T any](path string) []*fastjson.Value {
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

	var parser fastjson.Parser
	v, err := parser.ParseBytes(buf.Bytes())
	if err != nil {
		panic(err)
	}

	a, err := v.Array()
	if err != nil {
		panic(err)
	}

	return a
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

func getGamePath(gameID string) string {
	return fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID)
}

func getGame(gameID string) []*fastjson.Value {
	return parse[[]*fastjson.Value](fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID))
}

func analyze(events []*fastjson.Value) (*Game, error) {
	game := Game{
		Start: parseTime(string(events[0].GetStringBytes("eventTime"))),
		End:   parseTime(string(events[0].GetStringBytes("eventTime"))),
	}

	for _, event := range events {
		eventType := EventType(event.GetStringBytes("eventType"))
		switch eventType {
		case BULDING_DESTROYED:
		case CHAMPION_KILL:
		case CHAMPION_KILL_SPECIAL:
		case CHAMPION_LEVEL_UP:
		case EPIC_MONSTER_KILL:
		case GAME_INFO:
		case GAME_END:
		case ITEM_DESTROYED:
		case ITEM_PURCHASED:
		case ITEM_SOLD:
		case SKILL_LEVEL_UP:
		case STATS_UPDATE:
		case SUMMONER_SPELL_USED:
		case TURRET_PLATE_DESTROYED:
		case WARD_KILLED:
		case WARD_PLACED:
		case QUEUED_EPIC_MONSTER_INFO:
		case QUEUED_DRAGON_INFO:
		case EPIC_MONSTER_SPAWN:
		case TURRET_PLATE_GOLD_EARNED:
		case ITEM_UNDO:
		case OBJECTIVE_BOUNTY_PRESTART:
		case OBJECTIVE_BOUNTY_FINISH:
		case SURRENDER_VOTE_START:
		case SURRENDER_FAILED_VOTES:
		case SURRENDER_VOTE:
		case SURRENDER_AGREED:
		case CHAMPION_REVIVED:
		case CHAMPION_TRANSFORMED:
		case UNANIMOUS_SURRENDER_VOTE_START:
		case CHAMP_SELECT:
		default:
			panic(fmt.Sprintf("unknown event type: %s", eventType))
		}

		// fmt.Printf("event type %s\n", eventType)

		eventTime := parseTime(string(event.GetStringBytes("eventTime")))

		if eventTime.Before(game.Start) {
			game.Start = eventTime
		}
		if eventTime.After(game.End) {
			game.End = eventTime
		}
	}

	return &game, nil
}
