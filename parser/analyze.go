package main

import (
	"fmt"
	"time"

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
			stats_update(&game, obj, t)
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

func stats_update(game *Game, obj *simdjson.Object, now time.Time) {
	analyze_participants(game, obj, now)
	analyze_teams(game, obj, now)
}

func analyze_participants(game *Game, obj *simdjson.Object, now time.Time) {
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

		pid := obj.FindKey("participantID", nil)
		participantID, err := pid.Iter.Int()
		if err != nil {
			panic(err)
		}

		switch participantID {
		case 1: // red top
			analyze_participant(&game.Red.Top.Player, obj, now)
		case 2: // red jg
			analyze_participant(&game.Red.Jungle.Player, obj, now)
		case 3: // red mid
			analyze_participant(&game.Red.Mid.Player, obj, now)
		case 4: // red adc
			analyze_participant(&game.Red.Adc.Player, obj, now)
		case 5: // red support
			analyze_participant(&game.Red.Support.Player, obj, now)
		case 6: // blue top
			analyze_participant(&game.Blue.Top.Player, obj, now)
		case 7: // blue jg
			analyze_participant(&game.Blue.Jungle.Player, obj, now)
		case 8: // blue mid
			analyze_participant(&game.Blue.Mid.Player, obj, now)
		case 9: // blue adc
			analyze_participant(&game.Blue.Adc.Player, obj, now)
		case 10: // blue support
			analyze_participant(&game.Blue.Support.Player, obj, now)
		}

	}
}

func analyze_participant(player *Player, obj *simdjson.Object, now time.Time) {
	analyze_participant_stats(player, obj, now)

	x := obj.FindKey("XP", nil)
	xp, err := x.Iter.Int()
	if err != nil {
		panic(err)
	}

	player.XP.Set(int(xp), now)
}

func analyze_teams(game *Game, obj *simdjson.Object, now time.Time) {
	t := obj.FindKey("teams", nil)
	teams, err := t.Iter.Array(nil)
	if err != nil {
		panic(err)
	}

	iter := teams.Iter()
	for {
		typ := iter.Advance()
		if typ != simdjson.TypeObject {
			break
		}

		obj, err := iter.Object(nil)
		if err != nil {
			break
		}

		pid := obj.FindKey("teamID", nil)
		teamID, err := pid.Iter.Int()
		if err != nil {
			panic(err)
		}

		analyze := func(team *Team) {
			k := obj.FindKey("championsKills", nil)
			kills, err := k.Iter.Int()
			if err != nil {
				panic(err)
			}
			team.KDA.Kill.Set(int(kills), now)

			d := obj.FindKey("deaths", nil)
			deaths, err := d.Iter.Int()
			if err != nil {
				panic(err)
			}
			team.KDA.Death.Set(int(deaths), now)

			g := obj.FindKey("totalGold", nil)
			gold, err := g.Iter.Int()
			if err != nil {
				panic(err)
			}
			team.Gold.Set(int(gold), now)
		}
		switch teamID {

		case 100:
			analyze(&game.Red)
		case 200:
			analyze(&game.Blue)
		}
	}
}

func analyze_participant_stats(player *Player, obj *simdjson.Object, now time.Time) {
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
			v := obj.FindKey("value", nil)
			value, err := v.Iter.Int()
			if err != nil {
				panic(err)
			}
			player.KDA.Kill.Set(int(value), now)
		case "NUM_DEATHS":
			v := obj.FindKey("value", nil)
			value, err := v.Iter.Int()
			if err != nil {
				panic(err)
			}
			player.Death.Set(int(value), now)
		case "ASSISTS":
			v := obj.FindKey("value", nil)
			value, err := v.Iter.Int()
			if err != nil {
				panic(err)
			}
			player.Assist.Set(int(value), now)
		case "MINIONS_KILLED":
			v := obj.FindKey("value", nil)
			value, err := v.Iter.Int()
			if err != nil {
				panic(err)
			}
			player.CS.Set(int(value), now)
		case "TOTAL_DAMAGE_DEALT_TO_OBJECTIVES":
			v := obj.FindKey("value", nil)
			value, err := v.Iter.Float()
			if err != nil {
				panic(err)
			}
			player.ObjectiveDamage.Set(value, now)
		case "VISION_SCORE":
			v := obj.FindKey("value", nil)
			value, err := v.Iter.Int()
			if err != nil {
				panic(err)
			}
			player.VisionScore.Set(int(value), now)
		default:
			// not handled
		}
	}
}
