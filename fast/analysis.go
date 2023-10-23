package main

import (
	"encoding/json"
	"time"
)

type KDA struct {
	Kill   int `json:"kill"`
	Death  int `json:"death"`
	Assist int `json:"assist"`
}

type Ratio float32

type Role string

type Count int

type Duration time.Duration

type Player struct {
	Role `json:"role"`
	KDA
	KDARatio             Ratio `json:"kda_ratio"`
	VisionScore          Count `json:"vision_score"`
	Cs                   Count `json:"cs"`
	CsRatio              Ratio `json:"cs_ratio"`
	XP                   Count `json:"xp"`
	XPRatio              Ratio `json:"xp_ratio"`
	ObjectiveDamage      Count `json:"objective_damage"`
	ObjectiveDamageRatio Ratio `json:"objective_damage_ratio"`
	TurretPlateGold      Count `json:"turret_plate_gold"`
	TurretPlateGoldRatio Ratio `json:"turret_plate_gold_ratio"`
	TurretDestroyed      Count `json:"turret_destroyed"`
	TurretDestroyedRatio Ratio `json:"turret_destroyed_ratio"`
}

type Top struct {
}

type Mid struct {
}

type Jungle struct {
	Baron       Count `json:"baron"`
	Dragon      Count `json:"dragon"`
	BaronRatio  Ratio `json:"baron_ratio"`
	DragonRatio Ratio `json:"dragon_ratio"`
}

type Adc struct {
	FirstDeath Duration `json:"first_death"`
}

type Support struct {
}

type Side string // red or blue

type Team struct {
	KDA
	KDARatio Ratio `json:"kda_ratio"`
	Top      `json:"top"`
	Mid      `json:"mid"`
	Jungle   `json:"jungle"`
	Adc      `json:"adc"`
	Support  `json:"support"`
}

type Game struct {
	Red  Team `json:"red"`
	Blue Team `json:"blue"`

	Start time.Time `json:"start"`
	End   time.Time `json:"end"`

	FirstBlood           Side `json:"first_blood"`
	FirstTurretDestroyed Side `json:"first_turret_destroyed"`
}

func (g Game) String() string {
	b, _ := json.MarshalIndent(g, "", "    ")
	return string(b)
}
