package main

type KDA struct {
	Kill   int `json:"kill"`
	Death  int `json:"death"`
	Assist int `json:"assist"`
}

type Ratio float32

type Player struct {
	KDA
	KDARatio Ratio `json:"kda_ration"`
}

type Side string // red or blue

type Team struct {
	KDA
	KDARatio Ratio `json:"kda_ratio"`
}

type Game struct {
	FirstTurretDestoryed Side `json:"first_turret_destoryed"`
}
