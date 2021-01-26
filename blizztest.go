package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const urlOauth = "https://us.battle.net/oauth/token"
const urlSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?mana=%d&class=%s"

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
	id		string `json:"id"`
	secret	string `json:"secret"`
	Bearer	BearerData
}

func (auth ClientAuth) FromFile(secretfile string) {
	filed, err := ioutil.ReadFile(secretfile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(filed, auth)
}

/* testing function to print each card in a human way */
func (card *HearthCard) PrettyPrint() string {
	var bufout bytes.Buffer
	var cardform = `
    Name: %s
    Type: %d
    Rarity: %d
    Set: %d
    Class: %d
`
	bufout.WriteString(fmt.Sprintf(cardform, card.Name, card.TypeId, card.Rarity, card.SetId, card.Class))
	return bufout.String()
}

/* new stack o cards from a response body */
func (cards *HearthCards) New(body []byte) {
	err := json.Unmarshal(body, &cards)
	if err != nil {
		log.Fatal(err)
	}
}

/* list them all and print them for testing/debugging */
func (cards *HearthCards) List() {
	cardsDb := cards.All

	for _, ele := range(cardsDb) {
		fmt.Print(ele.PrettyPrint())
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
	req.SetBasicAuth(auth.id, auth.secret)
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
	auth.Bearer = bearer_ticket

	defer resp.Body.Close()
}

func GetCards() (cardDb HearthCards) {
    auth.Login()
	client := &http.Client{}
	req, err := http.NewRequest("GET", cardsurl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add(auth.BearerHeader())
    resp, err := client.Do(req)
    if err != nil {
    	log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
 	cardDb.New(body)
	return cardDb
}

/* test handler for our webapp good for lb health */
func testWebHandler(writer http.ResponseWriter, r *http.Request) {
}

func cardsWebHandler(writer http.ResponseWriter, r *http.Request) {
}

func indexWebHandler(writer http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/test", testWebHandler)
	http.HandleFunc("/cards.json", cardsWebHandler)
	http.HandleFunc("/index.html", indexWebHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}