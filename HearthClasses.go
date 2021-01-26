package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const urlHearthClasses = "https://us.api.blizzard.com/hearthstone/metadata/classes?locale=en_US"

type HearthClasses struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

func (classes *HearthClasses) Init(auth ClientAuth) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlHearthClasses, nil)
	req.Header.Add(auth.BearerHeader())
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(body))
	classes.NewFromBytes(body)
	defer resp.Body.Close()
}

func (classes *HearthClasses) NewFromBytes(body []byte) {
	err := json.Unmarshal(body, classes)
	if err != nil {
		log.Fatal(err)
	}
}

func (classes *HearthClasses) Enrich(db *HearthCards) {
}