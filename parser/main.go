package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	GAMES_DIRECTIORY = "../data/games"
	GAMES_FILE       = "../data/games.json"
)

func main() {
	games := parse[[]map[string]any](GAMES_FILE)

	for _, game := range games {
		gameID := game["id"].(string)

		data, err := analyze(getGame(gameID))
		if err != nil {
			panic(err)
		}

		fmt.Print(data)
	}
}

func parse[T any](path string) T {
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

	var data T
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

func getGame(gameID string) Events {
	return parse[Events](fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID))
}

func analyze(events Events) (*Game, error) {
	game := Game{
		Start: startTime(events),
		End:   endTime(events),
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
		default:
			panic(fmt.Sprintf("unknown event type: %s", t))
		}

		// fmt.Printf("event type %s\n", eventType)
	}

	return &game, nil
}

func startTime(events Events) time.Time {
	start := events[0].EventTime()

	for _, event := range events {
		t := event.EventTime()

		if t.Before(start) {
			start = event.EventTime()
		}
	}

	return start
}

func endTime(events Events) time.Time {
	start := events[0].EventTime()

	for _, event := range events {
		t := event.EventTime()

		if t.After(start) {
			start = event.EventTime()
		}
	}

	return start
}
