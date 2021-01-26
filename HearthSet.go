package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const urlHearthSets = "https://us.api.blizzard.com/hearthstone/metadata/sets?locale=en_US"

type HearthSet struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type HearthSets []HearthSet

func (sets *HearthSets) Init(auth ClientAuth) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlHearthSets, nil)
	req.Header.Add(auth.BearerHeader())
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(body))
	sets.NewFromBytes(body)
	defer resp.Body.Close()
}

/* new stack o cards from a response body */
func (sets *HearthSets) NewFromBytes(body []byte) {
	err := json.Unmarshal(body, sets)
	if err != nil {
		log.Fatal(err)
	}
}

func (sets *HearthSets) Enrich(db *HearthCards) {
}