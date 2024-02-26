package oauth_oidc

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/sso_client/pkg/state"
	"github.com/Chingizkhan/sso_client/service/sso_service_client"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

const (
	tokenIntrospectURL = "/oauth/token/introspect"
)

type (
	CookieProcessor interface {
		GenerateCookie(state string) http.Cookie
	}

	UseCase struct {
		addr         string
		cookie       CookieProcessor
		client       sso_service_client.Client
		oauth2Config oauth2.Config
	}
)

func New(addr string, cookie CookieProcessor, client sso_service_client.Client, oauth2Config oauth2.Config) *UseCase {
	return &UseCase{
		cookie:       cookie,
		addr:         addr,
		client:       client,
		oauth2Config: oauth2Config,
	}
}

func (u *UseCase) Login() (string, http.Cookie, error) {
	st, err := state.Generate()
	if err != nil {
		return "", http.Cookie{}, fmt.Errorf("state.Generate: %w", err)
	}

	cookie := u.cookie.GenerateCookie(string(st))
	loginUrl := u.oauth2Config.AuthCodeURL(string(st))
	return loginUrl, cookie, nil
}

func (u *UseCase) Introspect(ctx context.Context, accessToken string) error {
	introspect, err := u.client.Introspect(ctx, u.addr+tokenIntrospectURL, accessToken)
	if err != nil {
		return fmt.Errorf("OauthClient.Introspect: %w", err)
	}

	log.Println("oauth_oidc introspect response: ", introspect)

	return nil
}
