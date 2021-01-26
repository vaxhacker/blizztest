package main

import (
	"html/template"
	"log"
	"net/http"
)

const secretFile = "blizz-secret.json"
const urlOauth = "https://us.battle.net/oauth/token"
//const urlSearch = "https://us.api.blizzard.com/hearthstone/cards/?locale=en-US?mana=%d&class=%s"

/* test handler for our webapp good for lb health */
func testWebHandler(writer http.ResponseWriter, req *http.Request) {
}

func authAndFetchDb() HearthCardDb {
	auth := ClientAuth{}
	auth.FromWhatever() // stub for some type of secrets storage
	auth.Login()        // get bearer token
	db := InitCardDb(auth)       // initialize datastructures of cards
	log.Print(db)
	return db
}

func indexWebHandler(writer http.ResponseWriter, req *http.Request) {
	db := authAndFetchDb()
	tmpl := template.Must(template.ParseFiles("public_html/hearthstone.html"))
	tmpl.Execute(writer, db.AllSorted())
}

func main() {
	http.HandleFunc("/test", testWebHandler)
	http.HandleFunc("/hearthstone.html", indexWebHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}