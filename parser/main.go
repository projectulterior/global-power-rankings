package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/minio/simdjson-go"
)

const (
	GAMES_DIRECTIORY = "../data/games"
	GAMES_FILE       = "../data/games.json"
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

func analyze(events simdjson.ParsedJson) (*Game, error) {
	game := Game{}

	eventsIter := events.Iter()
	eventsIter.Advance()

	_, arrIter, err := eventsIter.Root(nil)
	if err != nil {
		return nil, err
	}

	gameArr, err := arrIter.Array(nil)
	if err != nil {
		return nil, err
	}

	gameArrIter := gameArr.Iter()
	for {
		typ := gameArrIter.Advance()
		if typ != simdjson.TypeObject {
			break
		}

		obj, err := gameArrIter.Object(nil)
		if err != nil {
			break
		}

		el := obj.FindKey("eventType", nil)
		eventType, err := el.Iter.String()
		if err != nil {
			return nil, err
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
			return nil, err
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

	return &game, nil
}
