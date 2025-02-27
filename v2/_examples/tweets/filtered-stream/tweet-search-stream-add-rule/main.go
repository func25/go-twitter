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

	In order to run, the user will need to provide the bearer token, rule, tag and run method.

*
*/
func main() {
	token := flag.String("token", "", "twitter API token")
	rule := flag.String("rule", "", "rules")
	tag := flag.String("tag", "", "tag")
	dryRun := flag.Bool("dry_run", false, "dry run")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to tweet search stream add rule callout")

	streamRule := twitter.TweetSearchStreamRule{
		Value: *rule,
		Tag:   *tag,
	}

	searchStreamRules, err := client.TweetSearchStreamAddRule(context.Background(), []twitter.TweetSearchStreamRule{streamRule}, *dryRun)
	if err != nil {
		log.Panicf("tweet search stream add rule callout error: %v", err)
	}

	enc, err := json.MarshalIndent(searchStreamRules, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

}
