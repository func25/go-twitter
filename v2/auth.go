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

type twOAuth1 struct {
	config *oauth1.Config
	token  *oauth1.Token
}

func (t *twOAuth1) Client(ctx context.Context) *http.Client {
	return t.config.Client(ctx, t.token)
}

func NewOAuth1Config(consumerKey, consumerSecret, accessToken, accessTokenSecret string) twOAuth1 {
	return twOAuth1{
		config: oauth1.NewConfig(consumerKey, consumerSecret),
		token:  oauth1.NewToken(accessToken, accessTokenSecret),
	}
}

// ==== Singleton
var globalTwiOAuth1 twOAuth1

func SetOAuth1(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	globalTwiOAuth1 = twOAuth1{
		config: oauth1.NewConfig(consumerKey, consumerSecret),
		token:  oauth1.NewToken(accessToken, accessTokenSecret),
	}
}

func NewOAuthClient(ctx context.Context) (*http.Client, error) {
	if globalTwiOAuth1.config == nil {
		return nil, ErrOAuth1NotSetup
	}

	return globalTwiOAuth1.Client(ctx), nil
}
