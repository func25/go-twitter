package twitter

import (
	"context"
	"errors"
	"net/http"

	"github.com/dghubble/oauth1"
)

var (
	ErrOAuth1NotSetup = errors.New("twitter: oauth1 not setup")
)

// Authorizer will add the authorization to the HTTP request
type Authorizer interface {
	Add(req *http.Request)
}

type TwOAuth1 struct {
	config *oauth1.Config
	token  *oauth1.Token
}

func (t *TwOAuth1) Client(ctx context.Context) *http.Client {
	return t.config.Client(ctx, t.token)
}

func (t *TwOAuth1) TwitterClient(ctx context.Context) *Client {
	return &Client{
		Client: t.Client(ctx),
		Host:   "https://api.twitter.com",
	}
}

func NewOAuth1Config(consumerKey, consumerSecret, accessToken, accessTokenSecret string) TwOAuth1 {
	return TwOAuth1{
		config: oauth1.NewConfig(consumerKey, consumerSecret),
		token:  oauth1.NewToken(accessToken, accessTokenSecret),
	}
}

// ==== Singleton
var globalTwiOAuth1 TwOAuth1

func SetOAuth1(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	globalTwiOAuth1 = TwOAuth1{
		config: oauth1.NewConfig(consumerKey, consumerSecret),
		token:  oauth1.NewToken(accessToken, accessTokenSecret),
	}
}

func NewOAuth1Client(ctx context.Context) (*http.Client, error) {
	if globalTwiOAuth1.config == nil {
		return nil, ErrOAuth1NotSetup
	}

	return globalTwiOAuth1.Client(ctx), nil
}

func NewTwitterV2OAuth1(ctx context.Context) (*Client, error) {
	if globalTwiOAuth1.config == nil {
		return nil, ErrOAuth1NotSetup
	}

	return globalTwiOAuth1.TwitterClient(ctx), nil
}
