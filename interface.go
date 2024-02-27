package sso_client

import (
	"context"
	"github.com/Chingizkhan/sso_client/service/sso_service_client"
	"golang.org/x/oauth2"
	"net/http"
)

type (
	ClientCredentials interface {
		Introspect(ctx context.Context, accessToken string) error
	}

	OauthOidc interface {
		Login() (string, http.Cookie, error)
		Introspect(ctx context.Context, accessToken string) (*sso_service_client.IntrospectResponse, error)
		Callback(ctx context.Context, code string) (*oauth2.Token, *sso_service_client.IntrospectResponse, error)
	}

	CookieProcessor interface {
		GenerateCookie(state string) http.Cookie
	}
)
