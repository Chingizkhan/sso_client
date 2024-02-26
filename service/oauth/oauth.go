package oauth

import (
	"golang.org/x/oauth2"
)

// todo: вытащить отдельно в конфиг
var (
	Config oauth2.Config
)

func init() {
	Config = oauth2.Config{
		ClientID:     "38b36b9d-48a8-40fd-9911-ee4462428c58",
		ClientSecret: "mysecret",
		RedirectURL:  "http://localhost:8082/callback",
		Scopes:       []string{"offline", "users.write", "users.read", "users.edit", "users.delete"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9010/oauth2/auth",
			TokenURL: "http://localhost:9010/oauth2/token",
		},
	}
}
