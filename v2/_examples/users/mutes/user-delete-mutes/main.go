package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	twitter "github.com/func25/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/*
*

	In order to run, the user will need to provide the bearer token and the list of tweet ids.

*
*/
func main() {
	token := flag.String("token", "", "twitter API token")
	userID := flag.String("user_id", "", "user id")
	targetUserID := flag.String("target_user_id", "", "target user id")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to delete user mute callout")

	userResponse, err := client.DeleteUserMutes(context.Background(), *userID, *targetUserID)
	if err != nil {
		log.Panicf("user delete mute error: %v", err)
	}

	enc, err := json.MarshalIndent(userResponse, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}
