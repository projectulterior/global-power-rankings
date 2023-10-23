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

	"github.com/minio/simdjson-go"
)

const (
	GAMES_DIRECTIORY = "../data/games"
	GAMES_FILE       = "../data/games.json"
)

func main() {
	run()
	// run_old()
}

func run_old() {
	games := parse[map[string]any](GAMES_FILE)
	fmt.Printf("length -- %d", len(games))

	var wg sync.WaitGroup
	sema := make(chan struct{}, 30)

	start := time.Now()

	var count atomic.Int32

	for i, g := range games {
		wg.Add(1)
		sema <- struct{}{}
		go func(game map[string]any) {
			defer wg.Done()
			defer func() { <-sema }()

			gameID := game["id"].(string)
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
				fmt.Println("no file found")
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

func getGamePath(gameID string) string {
	return fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID)
}

func getGame(gameID string) Events {
	return parse[Event](fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID))
}

func analyze(events Events) (*Game, error) {
	game := Game{
		Start: events[0].EventTime(),
		End:   events[0].EventTime(),
	}

	for i, event := range events {
		t, ok := event["eventType"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to parse event type: %d", i)
		}

		eventType := EventType(t)
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
			panic(fmt.Sprintf("unknown event type: %s", t))
		}

		// fmt.Printf("event type %s\n", eventType)

		eventTime := event.EventTime()

		if eventTime.Before(game.Start) {
			game.Start = eventTime
		}
		if eventTime.After(game.End) {
			game.End = eventTime
		}
	}

	return &game, nil
}

func simdAnal(events simdjson.ParsedJson) (*Game, error) {
	game := Game{}

	events.ForEach(func(i simdjson.Iter) error {
		switch i.Type() {
		case simdjson.TypeObject:
			obj, _ := i.Object(nil)
			el := obj.FindKey("eventType", nil)
			eventType, err := el.Iter.String()
			if err != nil {
				return err
			}

			switch EventType(eventType) {
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

			el = obj.FindKey("eventTime", nil)
			eventTime, err := el.Iter.String()
			if err != nil {
				return err
			}

			t := EventTimeSIMD(eventTime)
			if game.Start.IsZero() && game.End.IsZero() {
				game.Start = t
				game.End = t
			}
			if t.Before(game.Start) {
				game.Start = t
			}
			if t.After(game.End) {
				game.End = t
			}
		}

		return nil
	})

	return &game, nil
}
