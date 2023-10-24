package main

import "time"

type Event map[string]any
type Events []Event

type EventType string

const (
	BULDING_DESTROYED      EventType = "building_destroyed"
	CHAMPION_KILL          EventType = "champion_kill"
	CHAMPION_KILL_SPECIAL  EventType = "champion_kill_special"
	CHAMPION_LEVEL_UP      EventType = "champion_level_up"
	EPIC_MONSTER_KILL      EventType = "epic_monster_kill"
	GAME_INFO              EventType = "game_info"
	GAME_END               EventType = "game_end"
	ITEM_DESTROYED         EventType = "item_destroyed"
	ITEM_PURCHASED         EventType = "item_purchased"
	ITEM_SOLD              EventType = "item_sold"
	SKILL_LEVEL_UP         EventType = "skill_level_up"
	STATS_UPDATE           EventType = "stats_update"
	SUMMONER_SPELL_USED    EventType = "summoner_spell_used"
	TURRET_PLATE_DESTROYED EventType = "turret_plate_destroyed"
	WARD_KILLED            EventType = "ward_killed"
	WARD_PLACED            EventType = "ward_placed"

	// more events
	QUEUED_EPIC_MONSTER_INFO       EventType = "queued_epic_monster_info"
	QUEUED_DRAGON_INFO             EventType = "queued_dragon_info"
	EPIC_MONSTER_SPAWN             EventType = "epic_monster_spawn"
	TURRET_PLATE_GOLD_EARNED       EventType = "turret_plate_gold_earned"
	ITEM_UNDO                      EventType = "item_undo"
	OBJECTIVE_BOUNTY_PRESTART      EventType = "objective_bounty_prestart"
	OBJECTIVE_BOUNTY_FINISH        EventType = "objective_bounty_finish"
	SURRENDER_VOTE_START           EventType = "surrenderVoteStart"
	SURRENDER_FAILED_VOTES         EventType = "surrenderFailedVotes"
	SURRENDER_VOTE                 EventType = "surrenderVote"
	SURRENDER_AGREED               EventType = "surrenderAgreed"
	CHAMPION_REVIVED               EventType = "champion_revived"
	CHAMPION_TRANSFORMED           EventType = "champion_transformed"
	UNANIMOUS_SURRENDER_VOTE_START EventType = "unanimousSurrenderVoteStart"
	CHAMP_SELECT                   EventType = "champ_select"
)

func (e Event) EventTime() time.Time {
	s, ok := e["eventTime"].(string)
	if !ok {
		panic("error in parsing event time")
	}

	t, err := time.Parse("2006-01-02T15:04:05.999Z", s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.99Z", s)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05.9Z", s)
			if err != nil {
				t, err = time.Parse("2006-01-02T15:04:05Z", s)
				if err != nil {
					t, err = time.Parse("2006-01-02T15:04Z", s)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	return t
}

func EventTimeSIMD(eventTime string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.999Z", eventTime)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.99Z", eventTime)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05.9Z", eventTime)
			if err != nil {
				t, err = time.Parse("2006-01-02T15:04:05Z", eventTime)
				if err != nil {
					t, err = time.Parse("2006-01-02T15:04Z", eventTime)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	return t
}
