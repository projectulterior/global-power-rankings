package main

import (
	"fmt"

	"github.com/minio/simdjson-go"
)

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

		el := obj.FindKey("eventTime", nil)
		eventTime, err := el.Iter.String()
		if err != nil {
			return nil, err
		}

		// event timestamp
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

		// event type
		el = obj.FindKey("eventType", nil)
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
			stats_update(&game, obj)
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

	}

	return &game, nil
}

func stats_update(game *Game, obj *simdjson.Object) {
	p := obj.FindKey("participants", nil)

	participants, err := p.Iter.Array(nil)
	if err != nil {
		panic(err)
	}

	iter := participants.Iter()
	for {
		typ := iter.Advance()
		if typ != simdjson.TypeObject {
			break
		}

		obj, err := iter.Object(nil)
		if err != nil {
			break
		}

		e := obj.FindKey("participantID", nil)
		participantID, err := e.Iter.Int()
		if err != nil {
			panic(err)
		}
		fmt.Println(participantID)

	}
}

func analyze_stats(game *Game, obj *simdjson.Object) {
	s := obj.FindKey("stats", nil)
	stats, err := s.Iter.Array(nil)
	if err != nil {
		panic(err)
	}

	iter := stats.Iter()
	for {
		typ := iter.Advance()
		if typ != simdjson.TypeObject {
			break
		}

		obj, err := iter.Object(nil)
		if err != nil {
			break
		}

		n := obj.FindKey("name", nil)
		name, err := n.Iter.String()
		if err != nil {
			panic(err)
		}

		switch name {
		case "CHAMPIONS_KILLED":
			v := obj.FindKey("vale", nil)
			value, err := v.Iter.Int()
			if err != nil {
				panic(err)
			}

			fmt.Println("champion killed", value)
		}
	}
}
