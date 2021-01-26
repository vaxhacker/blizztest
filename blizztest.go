package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const secretFile = "blizz-secret.json"
const urlOauth = "https://us.battle.net/oauth/token"
//const urlSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?mana=%d&class=%s"
const urlSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?class=warlock"

type HearthCard struct {
	Id 			int	`json:"id"`
	Image		string `json:"image"`
	Name		string `json:"name"`
	TypeId		int `json:"cardTypeId"`
	SetId		int `json:"cardSetId"`
	Rarity		int `json:"rarityId"`
	Class 		int `json:"classId"`
}

type HearthCards struct {
	All			[]HearthCard	`json:"cards"`
}

type BearerData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}


type ClientAuth struct {
	Id		string `json:"id"`
	Secret	string `json:"secret"`
	Bearer	BearerData
}

func (auth *ClientAuth) FromWhatever() {
	data, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, &auth)
	log.Print(auth)
}

func (db *HearthCards) Init(auth ClientAuth) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlSearch, nil)
	req.Header.Add(auth.BearerHeader())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	log.Print(string(body))
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

/* return a bearer header to be used with httpclient */
func (auth *ClientAuth) BearerHeader() (string, string) {
	return "Authorization", fmt.Sprintf("Bearer %s", auth.Bearer.AccessToken)
}

/* create a new clientauth struct takes id and secret */
func (auth *ClientAuth) Login() {
	data := url.Values{}
	data.Add("grant_type","client_credentials")

	/* build a client we will use to request our bearer token.
	   we cannot use the httpClient.Post methods as we are required to
	   set basic auth on the request object.
	 */
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlOauth, bytes.NewBufferString(data.Encode()))
	req.SetBasicAuth(auth.Id, auth.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.PostForm = data
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	/* read our response into body */
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	/* deserialize the json and store the bearer in BearerData struct */
	bearer := BearerData{}
	err = json.Unmarshal(body, &bearer)
	if err != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), err.Error())
	}
	auth.Bearer = bearer

	defer resp.Body.Close()
}

/* test handler for our webapp good for lb health */
func testWebHandler(writer http.ResponseWriter, req *http.Request) {
}

func cardsWebHandler(writer http.ResponseWriter, req *http.Request) {
	auth := ClientAuth{}
	auth.FromWhatever()
	log.Print(auth)
	auth.Login()
	log.Print(auth)

	db := HearthCards{}
	db.Init(auth)
}

func indexWebHandler(writer http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("public_html/index.html"))
	tmpl.Execute(writer, nil)
}

func main() {
	http.HandleFunc("/test", testWebHandler)
	http.HandleFunc("/cards.json", cardsWebHandler)
	http.HandleFunc("/hearthstone.html", indexWebHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}