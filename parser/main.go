package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	GAMES_DIRECTIORY = "../data/games/"
	GAMES_FILE       = "../data/games.json"
)

func main() {
	gameID := "110733838936446929"

	analyze(getGame(gameID))
}

func getGame(gameID string) Events {
	file, err := os.Open(fmt.Sprintf("%s/%s.json", GAMES_DIRECTIORY, gameID))
	if err != nil {
		panic(err)
	}

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

func analyze(events Events) error {
	for i, event := range events {
		t, ok := event["eventType"].(string)
		if !ok {
			return fmt.Errorf("unable to parse event type: %d", i)
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
		default:
			panic(fmt.Sprintf("unknown event type: %s", t))
		}

		fmt.Printf("event type %s\n", eventType)
	}

	return nil
}
