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

	In order to run, the user will need to provide the bearer token and a query.

*
*/
func main() {
	token := flag.String("token", "", "twitter API token")
	query := flag.String("query", "", "twitter query")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.TweetAllCountsOpts{
		Granularity: twitter.GranularityHour,
	}

	fmt.Println("Callout to tweet all counts callout")

	tweetResponse, err := client.TweetAllCounts(context.Background(), *query, opts)
	if err != nil {
		log.Panicf("tweet all counts error: %v", err)
	}

	enc, err := json.MarshalIndent(tweetResponse.TweetCounts, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

	metaBytes, err := json.MarshalIndent(tweetResponse.Meta, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(metaBytes))
}
