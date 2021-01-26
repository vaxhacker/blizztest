package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

const urlCardSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?class=%s&manaCost=%d"
const minMana = 7
const maxMana = 10
var classesToSearch = []string { "warlock", "druid" }

type HearthCard struct {
	Id 			int	`json:"id"`
	Image		string `json:"image"`
	Name		string `json:"name"`
	FancyName   string
	TypeId		int `json:"cardTypeId"`
	FancyType   string
	SetId		int `json:"cardSetId"`
	FancySet    string
	Rarity		int `json:"rarityId"`
	Class 		int `json:"classId"`
	FancyClass  string
	ManaCost    int `json:"manaCost"`
}

type HearthCards struct {
	All			[]HearthCard	`json:"cards"`
}

func (db HearthCards) AllSorted() []HearthCard {
	sorted := db.All
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Id < sorted[j].Id })
	return sorted
}

func (db *HearthCards) ToJson() string {
	data, err := json.Marshal(db)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (db *HearthCards) GetAllCards(auth ClientAuth) {
	outputCards := make([]HearthCard, 0, 10)
	client := &http.Client{}
	for _, ele := range classesToSearch {
		for i := minMana; i < maxMana+1; i++ {
			shim := new(HearthCards)
			search := fmt.Sprintf(urlCardSearch, ele, i)
			req, _ := http.NewRequest("GET", search, nil)
			req.Header.Add(auth.BearerHeader())
			resp, _ := client.Do(req)
			body, _ := ioutil.ReadAll(resp.Body) //need to close this.
			err := json.Unmarshal([]byte(body), shim)
			if err != nil {
				log.Print("sad")
				log.Fatal(err)
			}
			for _, ele :=  range shim.All {
				outputCards = append(outputCards, ele)
			}
		}
	}
	db.All = outputCards
}

func (db *HearthCards) Init(auth ClientAuth) {
	db.GetAllCards(auth)
	classes := HearthClassDb{}
	sets := HearthSetDb{}
	classes.Init(auth, db)
	sets.Init(auth, db)
}
