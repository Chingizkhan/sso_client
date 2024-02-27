package oauth_oidc

import (
	"context"
	"errors"
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
		cookie       CookieProcessor
		client       sso_service_client.Client
		oauth2Config oauth2.Config
	}
)

func New(cookie CookieProcessor, client sso_service_client.Client, oauth2Config oauth2.Config) *UseCase {
	return &UseCase{
		cookie:       cookie,
		client:       client,
		oauth2Config: oauth2Config,
	}
}

func (u *UseCase) Login() (string, http.Cookie, error) {
	st, err := state.Generate()
	if err != nil {
		return "", http.Cookie{}, fmt.Errorf("state.Generate: %w", err)
	}
	log.Println("state:", string(st))

	cookie := u.cookie.GenerateCookie(string(st))
	loginUrl := u.oauth2Config.AuthCodeURL(string(st))
	return loginUrl, cookie, nil
}

func (u *UseCase) Introspect(ctx context.Context, accessToken string) (*sso_service_client.IntrospectResponse, error) {
	introspect, err := u.client.Introspect(ctx, tokenIntrospectURL, accessToken)
	if err != nil {
		return nil, fmt.Errorf("OauthClient.Introspect: %w", err)
	}

	log.Println("oauth_oidc introspect response: ", introspect)

	return introspect, nil
}

func (u *UseCase) Callback(ctx context.Context, code string) (*oauth2.Token, *sso_service_client.IntrospectResponse, error) {
	token, err := u.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, nil, errors.New("can not exchange token: " + err.Error())
	}

	introspectResponse, err := u.Introspect(ctx, token.AccessToken)
	if err != nil {
		return nil, nil, errors.New("can not introspect response: " + err.Error())
	}

	log.Println("client_id", introspectResponse.ClientID)
	log.Println("active", introspectResponse.Active)
	return token, introspectResponse, nil
}
