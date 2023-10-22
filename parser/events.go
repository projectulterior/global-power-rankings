package main

type Events []map[string]any

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
	QUEUED_EPIC_MONSTER_INFO  EventType = "queued_epic_monster_info"
	QUEUED_DRAGON_INFO        EventType = "queued_dragon_info"
	EPIC_MONSTER_SPAWN        EventType = "epic_monster_spawn"
	TURRET_PLATE_GOLD_EARNED  EventType = "turret_plate_gold_earned"
	ITEM_UNDO                 EventType = "item_undo"
	OBJECTIVE_BOUNTY_PRESTART EventType = "objective_bounty_prestart"
	OBJECTIVE_BOUNTY_FINISH   EventType = "objective_bounty_finish"
)
