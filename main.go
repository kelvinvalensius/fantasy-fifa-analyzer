package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

type Player struct {
	PlayerName         string
	CountryCode        string
	Skill              int
	OverallPoint       float64
	SelectedPercentage float64
	PlayerValue        float64
	Score              float64 `json:"-"`
}

type Data struct {
	Value []Player
}

type Response struct {
	Data Data
}

const (
	SkillGoalkeeper int = 1
	SkillDefender   int = 2
	SkillMidfield   int = 3
	SkillForward    int = 4
)

func printPlayersToCsv(players []Player, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	f.WriteString("PlayerName,CountryCode,Skill,OverallPoint,SelectedPercentage,PlayerValue,Score\n")
	for _, p := range players {
		s := fmt.Sprintf("%s,%s,%d,%f,%f,%f,%f\n", p.PlayerName, p.CountryCode, p.Skill, p.OverallPoint, p.SelectedPercentage, p.PlayerValue, p.Score)
		f.WriteString(s)
	}

	f.Sync()
}

func main() {
	resp, err := http.Get("https://fantasy.fifa.com/services/api/statistics/statsdetail?optType=1&gamedayId=1&language=en&buster=default")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := &Response{}
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Fatal(err)
		return
	}

	players := r.Data.Value
	for i, p := range players {
		players[i].Score = p.OverallPoint / p.PlayerValue * p.SelectedPercentage
	}

	sort.Slice(players, func(i, j int) bool {
		if players[i].Score != players[j].Score {
			return players[i].Score > players[j].Score
		}
		return players[i].PlayerValue < players[j].PlayerValue
	})

	printPlayersToCsv(players, "list.csv")
}
