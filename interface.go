package sso_client

import (
	"context"
	"net/http"
)

type (
	ClientCredentials interface {
		Introspect(ctx context.Context, accessToken string) error
	}

	OauthOidc interface {
		Login() (string, http.Cookie, error)
		Introspect(ctx context.Context, accessToken string) error
	}

	CookieProcessor interface {
		GenerateCookie(state string) http.Cookie
	}
)
