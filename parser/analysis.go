package main

import (
	"encoding/json"
	"time"
)

type KDA struct {
	Kill   CountInt `json:"kill"`
	Death  CountInt `json:"death"`
	Assist CountInt `json:"assist"`
}

type Ratio struct {
	last      time.Time
	durations map[Side]time.Duration
}

func (r *Ratio) Add(side Side, now time.Time) {
	dur := now.Sub(r.last)
	r.durations[side] += dur
	r.last = now
}

func (r *Ratio) Get(side Side) float32 {
	var this float32
	var total float32
	for s, dur := range r.durations {
		if s == side {
			this = float32(dur.Milliseconds())
		}

		total += float32(dur.Milliseconds())
	}

	return this / total
}

type Role string

type CountInt = Count[int]
type CountFloat = Count[float64]

type Count[T int | float64] struct {
	count T
	last  time.Time
}

func (c *Count[T]) Set(value T, now time.Time) {
	if now.After(c.last) {
		c.count = value
	}
}

func (c *Count[T]) Get() T {
	return c.count
}

type Duration time.Duration

type Player struct {
	Champion string `json:"champion"`
	Role     `json:"role"`
	KDA
	KDARatio               Ratio      `json:"kda_ratio"`
	VisionScore            CountInt   `json:"vision_score"`
	CS                     CountInt   `json:"cs"`
	CSRatio                Ratio      `json:"cs_ratio"`
	XP                     CountInt   `json:"xp"`
	XPRatio                Ratio      `json:"xp_ratio"`
	DamageToChampions      CountFloat `json:"damage_to_champions"`
	DamageToChampionsRatio Ratio      `json:"damage_to_champions_ratio"`
	ObjectiveDamage        CountFloat `json:"objective_damage"`
	ObjectiveDamageRatio   Ratio      `json:"objective_damage_ratio"`
	TurretPlateGold        CountInt   `json:"turret_plate_gold"`
	TurretPlateGoldRatio   Ratio      `json:"turret_plate_gold_ratio"`
	TurretDestroyed        CountInt   `json:"turret_destroyed"`
	TurretDestroyedRatio   Ratio      `json:"turret_destroyed_ratio"`
	MaxLevel               Duration   `json:"max_level"`
}

type Top struct {
	Player
}

type Mid struct {
	Player
}

type Jungle struct {
	Player
	Baron       CountInt `json:"baron"`
	Dragon      CountInt `json:"dragon"`
	BaronRatio  Ratio    `json:"baron_ratio"`
	DragonRatio Ratio    `json:"dragon_ratio"`
}

type Adc struct {
	Player
	FirstDeath Duration `json:"first_death"`
}

type Support struct {
	Player
}

type Side string // red or blue
const (
	RED  Side = "red"
	BLUE Side = "blue"
)

type Team struct {
	KDA
	KDARatio  Ratio    `json:"kda_ratio"`
	Gold      CountInt `json:"gold"`
	GoldRatio Ratio    `json:"gold_ratio"`
	Top       `json:"top"`
	Mid       `json:"mid"`
	Jungle    `json:"jungle"`
	Adc       `json:"adc"`
	Support   `json:"support"`
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
