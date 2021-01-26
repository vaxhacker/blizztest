package main

import (
	"html/template"
	"log"
	"net/http"
)

const secretFile = "blizz-secret.json"
const urlOauth = "https://us.battle.net/oauth/token"
//const urlSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?mana=%d&class=%s"
const urlSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?class=warlock"

/* test handler for our webapp good for lb health */
func testWebHandler(writer http.ResponseWriter, req *http.Request) {
}

func indexWebHandler(writer http.ResponseWriter, req *http.Request) {
	var db = HearthCards{}
	var auth = ClientAuth{}

	auth.FromWhatever()
	auth.Login()
	db.Init(auth)

	log.Print(db.AllSorted())
	tmpl := template.Must(template.ParseFiles("public_html/hearthstone.html"))
	tmpl.Execute(writer, []HearthCard(db.AllSorted()))
}

func main() {
	http.HandleFunc("/test", testWebHandler)
	http.HandleFunc("/hearthstone.html", indexWebHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}