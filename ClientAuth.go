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