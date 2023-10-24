package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	simd "github.com/minio/simdjson-go"
)

const WORKERS = 30

func run() {
	games := parse[map[string]any](GAMES_FILE)
	mapping := parse[map[string]any](MAPPING_FILE)

	var wg sync.WaitGroup
	sema := make(chan struct{}, WORKERS)

	var count atomic.Int32

	start := time.Now()
	for i, g := range games {
		wg.Add(1)
		sema <- struct{}{}
		go func(index int, game map[string]any, mapping []map[string]any) {
			defer wg.Done()
			defer func() { <-sema }()

			gameID := game["id"].(string)

			parsed, err := simdParse(getGamePath(gameID))
			if parsed == nil {
				fmt.Printf("parsed is null -- %s\n", err)
				return
			}

			iter := parsed.Iter()
			if iter.PeekNext() == simd.TypeNone {
				fmt.Println("no file found")
				return
			}

			analysis, err := analyze(*parsed)
			if err != nil {
				fmt.Printf("analysis error -- %s\n", err)
			}

			game["start_time"] = analysis.Start
			game["end_time"] = analysis.End

			// red
			game["red_kills"] = analysis.Red.KDA.Kill.Get()
			game["red_deaths"] = analysis.Red.KDA.Death.Get()
			game["red_gold"] = analysis.Red.Gold.Get()

			game["red_top_kills"] = analysis.Red.Top.KDA.Kill.Get()
			game["red_top_deaths"] = analysis.Red.Top.KDA.Death.Get()
			game["red_top_assists"] = analysis.Red.Top.KDA.Assist.Get()
			game["red_jungle_kills"] = analysis.Red.Jungle.KDA.Kill.Get()
			game["red_jungle_deaths"] = analysis.Red.Jungle.KDA.Death.Get()
			game["red_jungle_assists"] = analysis.Red.Jungle.KDA.Assist.Get()
			game["red_mid_kills"] = analysis.Red.Mid.KDA.Kill.Get()
			game["red_mid_deaths"] = analysis.Red.Mid.KDA.Death.Get()
			game["red_mid_assists"] = analysis.Red.Mid.KDA.Assist.Get()
			game["red_adc_kills"] = analysis.Red.Adc.KDA.Kill.Get()
			game["red_adc_deaths"] = analysis.Red.Adc.KDA.Death.Get()
			game["red_adc_assists"] = analysis.Red.Adc.KDA.Assist.Get()
			game["red_support_kills"] = analysis.Red.Support.KDA.Kill.Get()
			game["red_support_deaths"] = analysis.Red.Support.KDA.Death.Get()
			game["red_support_assists"] = analysis.Red.Support.KDA.Assist.Get()

			game["red_top_minion_kills"] = analysis.Red.Top.CS.Get()
			game["red_jungle_minion_kills"] = analysis.Red.Jungle.CS.Get()
			game["red_mid_minion_kills"] = analysis.Red.Mid.CS.Get()
			game["red_adc_minion_kills"] = analysis.Red.Adc.CS.Get()
			game["red_support_minion_kills"] = analysis.Red.Support.CS.Get()

			game["red_top_damage_to_objectives"] = analysis.Red.Top.ObjectiveDamage.Get()
			game["red_jungle_damage_to_objectives"] = analysis.Red.Jungle.ObjectiveDamage.Get()
			game["red_mid_damage_to_objectives"] = analysis.Red.Mid.ObjectiveDamage.Get()
			game["red_adc_damage_to_objectives"] = analysis.Red.Adc.ObjectiveDamage.Get()
			game["red_support_damage_to_objectives"] = analysis.Red.Support.ObjectiveDamage.Get()

			game["red_top_vision_score"] = analysis.Red.Top.VisionScore.Get()
			game["red_jungle_vision_score"] = analysis.Red.Jungle.VisionScore.Get()
			game["red_mid_vision_score"] = analysis.Red.Mid.VisionScore.Get()
			game["red_adc_vision_score"] = analysis.Red.Adc.VisionScore.Get()
			game["red_support_vision_score"] = analysis.Red.Support.VisionScore.Get()

			game["red_top_xp"] = analysis.Red.Top.XP.Get()
			game["red_jungle_xp"] = analysis.Red.Jungle.XP.Get()
			game["red_mid_xp"] = analysis.Red.Mid.XP.Get()
			game["red_adc_xp"] = analysis.Red.Adc.XP.Get()
			game["red_support_xp"] = analysis.Red.Support.XP.Get()

			// blue
			game["blue_kills"] = analysis.Blue.KDA.Kill.Get()
			game["blue_deaths"] = analysis.Blue.KDA.Death.Get()
			game["blue_gold"] = analysis.Blue.Gold.Get()

			game["blue_top_kills"] = analysis.Blue.Top.KDA.Kill.Get()
			game["blue_top_deaths"] = analysis.Blue.Top.KDA.Death.Get()
			game["blue_top_assists"] = analysis.Blue.Top.KDA.Assist.Get()
			game["blue_jungle_kills"] = analysis.Blue.Jungle.KDA.Kill.Get()
			game["blue_jungle_deaths"] = analysis.Blue.Jungle.KDA.Death.Get()
			game["blue_jungle_assists"] = analysis.Blue.Jungle.KDA.Assist.Get()
			game["blue_mid_kills"] = analysis.Blue.Mid.KDA.Kill.Get()
			game["blue_mid_deaths"] = analysis.Blue.Mid.KDA.Death.Get()
			game["blue_mid_assists"] = analysis.Blue.Mid.KDA.Assist.Get()
			game["blue_adc_kills"] = analysis.Blue.Adc.KDA.Kill.Get()
			game["blue_adc_deaths"] = analysis.Blue.Adc.KDA.Death.Get()
			game["blue_adc_assists"] = analysis.Blue.Adc.KDA.Assist.Get()
			game["blue_support_kills"] = analysis.Blue.Support.KDA.Kill.Get()
			game["blue_support_deaths"] = analysis.Blue.Support.KDA.Death.Get()
			game["blue_support_assists"] = analysis.Blue.Support.KDA.Assist.Get()

			game["blue_top_minion_kills"] = analysis.Blue.Top.CS.Get()
			game["blue_jungle_minion_kills"] = analysis.Blue.Jungle.CS.Get()
			game["blue_mid_minion_kills"] = analysis.Blue.Mid.CS.Get()
			game["blue_adc_minion_kills"] = analysis.Blue.Adc.CS.Get()
			game["blue_support_minion_kills"] = analysis.Blue.Support.CS.Get()

			game["blue_top_damage_to_objectives"] = analysis.Blue.Top.ObjectiveDamage.Get()
			game["blue_jungle_damage_to_objectives"] = analysis.Blue.Jungle.ObjectiveDamage.Get()
			game["blue_mid_damage_to_objectives"] = analysis.Blue.Mid.ObjectiveDamage.Get()
			game["blue_adc_damage_to_objectives"] = analysis.Blue.Adc.ObjectiveDamage.Get()
			game["blue_support_damage_to_objectives"] = analysis.Blue.Support.ObjectiveDamage.Get()

			game["blue_top_vision_score"] = analysis.Blue.Top.VisionScore.Get()
			game["blue_jungle_vision_score"] = analysis.Blue.Jungle.VisionScore.Get()
			game["blue_mid_vision_score"] = analysis.Blue.Mid.VisionScore.Get()
			game["blue_adc_vision_score"] = analysis.Blue.Adc.VisionScore.Get()
			game["blue_support_vision_score"] = analysis.Blue.Support.VisionScore.Get()

			game["red_top_xp"] = analysis.Red.Top.XP.Get()
			game["red_jungle_xp"] = analysis.Red.Jungle.XP.Get()
			game["red_mid_xp"] = analysis.Red.Mid.XP.Get()
			game["red_adc_xp"] = analysis.Red.Adc.XP.Get()
			game["red_support_xp"] = analysis.Red.Support.XP.Get()

			games[index] = game
		}(i, g, mapping)

		if i%100 == 0 {
			fmt.Printf("checkpoint -- %d -- %s\n", i, time.Since(start))
		}
	}
	wg.Wait()
	// close(data)

	b, err := json.MarshalIndent(games, "", "    ")
	if err != nil {
		panic(err)
	}

	write(OUTPUT_PATH, bytes.NewBuffer(b))
	fmt.Printf("errors: %d -- %s\n", count.Load(), time.Since(start))
	fmt.Println("------------")
}

func simdParse(path string) (*simd.ParsedJson, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	parsed, err := simd.Parse(buf.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
