package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

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

func (db *HearthCards) Init(auth ClientAuth) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlSearch, nil)
	req.Header.Add(auth.BearerHeader())
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	db.NewFromBytes(body)
	defer resp.Body.Close()
}

/* new stack o cards from a response body */
func (db *HearthCards) NewFromBytes(body []byte) {
	err := json.Unmarshal(body, &db)
	if err != nil {
		log.Fatal(err)
	}
}