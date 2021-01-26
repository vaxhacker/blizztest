package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

/* HearthClass, HeathSet, and any other metadata should probably moved to a generic
   metadata struct and support functions rather than the duplicate code here.
 */

const urlHearthClasses = "https://us.api.blizzard.com/hearthstone/metadata/classes?locale=en_US"
const urlHearthSets = "https://us.api.blizzard.com/hearthstone/metadata/sets?locale=en_US"

type HearthSet struct {
	Id HearthSetId `json:"id"`
	Name string `json:"name"`
}

type HearthSets []HearthSet
type HearthSetId int
type HearthSetDb map[HearthSetId]string

type HearthClass struct {
	Id HearthClassId `json:"id"`
	Name string   	 `json:"name"`
}
type HearthClasses []HearthClass
type HearthClassId int
type HearthClassDb map[HearthClassId]string

func (inv HearthSetDb) Init(auth ClientAuth, cards *HearthCards) {
	raw := HearthSets{}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlHearthSets, nil)
	req.Header.Add(auth.BearerHeader())
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &raw)
	if err != nil {
		log.Fatal(err)
	}
	for _, ele := range raw {
		inv[ele.Id] = ele.Name
	}
	inv.Enrich(cards)
	defer resp.Body.Close()
}

func (inv HearthClassDb) Init(auth ClientAuth, cards *HearthCards) {
	raw := HearthClasses{}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlHearthClasses, nil)
	req.Header.Add(auth.BearerHeader())
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &raw)
	if err != nil {
		log.Fatal(err)
	}
	for _, ele := range raw {
		inv[ele.Id] = ele.Name
	}
	inv.Enrich(cards)
	defer resp.Body.Close()
}

func (inv HearthSetDb) Enrich(db *HearthCards)  {
	dbptr := db.All
	for idx, ele := range dbptr {
		dbptr[idx].FancySet = inv[HearthSetId(ele.SetId)]
	}
}

func (inv HearthClassDb) Enrich(db *HearthCards)  {
	dbptr := db.All
	for idx, ele := range dbptr {
		dbptr[idx].FancyClass = inv[HearthClassId(ele.Class)]
	}
}