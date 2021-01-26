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

type HearthCardDb map[int]HearthCard

/* Sort the keys of the hash so we can get it in proper order */
func (db HearthCardDb) AllSorted() []HearthCard {
	sorted := make([]HearthCard, 0, len(db))
	keys := make([]int, 0, len(db))
	for key := range db {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		sorted = append(sorted, db[key])
	}
	return sorted
}

func (db *HearthCardDb) GetAllCards(auth ClientAuth) {
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
				log.Fatal(err)
			}
			for _, ele :=  range shim.All {
				(*db)[ele.Id] = ele
			}
		}
	}
}

func InitCardDb(auth ClientAuth) HearthCardDb {
	db := HearthCardDb{}
	classes := HearthClassDb{}
	sets := HearthSetDb{}

	db.GetAllCards(auth)
	classes.Init(auth, db)
	sets.Init(auth, db)

	return db
}
