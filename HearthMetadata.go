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

type HearthMeta struct {
	Id HearthId `json:"id"`
	Name string `json:"name"`
}

type HearthId int
type HearthClass HearthMeta
type HearthSet HearthMeta
type HearthSets []HearthSet
type HearthClasses []HearthClass
type HearthClassDb map[HearthId]string
type HearthSetDb map[HearthId]string

/* i want to use the HearthMeta and make a generic function
   that works on both sets and classes because they share id/name
   would reduce the repeated code!
 */
func (inv HearthSetDb) Init(auth ClientAuth, db HearthCardDb) {
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
	inv.Enrich(db)
	defer resp.Body.Close()
}

func (inv HearthClassDb) Init(auth ClientAuth, db HearthCardDb) {
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
	inv.Enrich(db)
	defer resp.Body.Close()
}

/* these would also be reduced with the HearthMeta
   type.
 */
func (inv HearthSetDb) Enrich(db HearthCardDb)  {
	for key, ele := range db {
		ele.FancySet = inv[HearthId(ele.SetId)]
		db[key] = ele
	}
}

func (inv HearthClassDb) Enrich(db HearthCardDb)  {
	for key, ele := range db {
		ele.FancyClass = inv[HearthId(ele.Class)]
		db[key] = ele
	}
}