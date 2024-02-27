package sso_client

import (
	"errors"
	"fmt"
	"github.com/Chingizkhan/sso_client/pkg/api_util"
	"github.com/Chingizkhan/sso_client/pkg/cookies"
	"github.com/Chingizkhan/sso_client/pkg/state"
	"github.com/Chingizkhan/sso_client/pkg/token"
	"github.com/Chingizkhan/sso_client/service/cookie_processor"
	"github.com/Chingizkhan/sso_client/service/sso_service_client"
	"github.com/Chingizkhan/sso_client/usecase/client_credentials"
	"github.com/Chingizkhan/sso_client/usecase/oauth_oidc"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type (
	Client interface {
		AuthOauth2(next http.Handler) http.Handler
		AuthClientCredentials(next http.Handler) http.Handler
		ProcessCallback(r *http.Request) (*oauth2.Token, *sso_service_client.IntrospectResponse, error)
	}

	SsoClient struct {
		cfg               Config
		oidc              OauthOidc
		clientCredentials ClientCredentials
		cookieProcessor   CookieProcessor
	}

	Config struct {
		CookieSecret   string
		CookieLifetime time.Duration
		OauthAddr      string
		Oauth2Config   oauth2.Config
	}
)

func New(config Config) *SsoClient {
	cookieProcessor := cookie_processor.New(config.CookieLifetime)
	ssoServiceClient := sso_service_client.New(time.Second*15, config.OauthAddr)
	oidc := oauth_oidc.New(cookieProcessor, ssoServiceClient, config.Oauth2Config)
	clientCredentials := client_credentials.New(ssoServiceClient)

	return &SsoClient{
		cfg:               config,
		oidc:              oidc,
		clientCredentials: clientCredentials,
		cookieProcessor:   cookieProcessor,
	}
}

func (c *SsoClient) AuthOauth2(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// check token exists
		accessToken := token.FromHttpRequest(r)
		if accessToken == "" {
			// get cookie and login url
			loginURL, cookie, err := c.oidc.Login()
			if err != nil {
				api_util.RenderErrorResponse(w, "can not show login flow", http.StatusInternalServerError)
				return
			}
			// set cookie to response
			err = cookies.WriteSigned(w, cookie, []byte(c.cfg.CookieSecret))
			if err != nil {
				api_util.RenderErrorResponse(w, "can not write signed cookies", http.StatusInternalServerError)
				return
			}
			// redirect
			http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
			return
		}

		// introspect token
		introspectResponse, err := c.oidc.Introspect(r.Context(), accessToken)
		if err != nil {
			api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !introspectResponse.Active {
			api_util.RenderErrorResponse(w, "token is inactive", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (c *SsoClient) AuthClientCredentials(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// check token exists
		accessToken := token.FromHttpRequest(r)
		if accessToken == "" {
			api_util.RenderErrorResponse(w, "empty token", http.StatusUnauthorized)
			return
		}

		// introspect token
		introspectResponse, err := c.clientCredentials.Introspect(r.Context(), accessToken)
		if err != nil {
			api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !introspectResponse.Active {
			api_util.RenderErrorResponse(w, "token is inactive", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (c *SsoClient) ProcessCallback(r *http.Request) (*oauth2.Token, *sso_service_client.IntrospectResponse, error) {
	var req CallbackRequest
	err := req.Validate(r)
	if err != nil {
		return nil, nil, fmt.Errorf("req.Validate: %w", err)
	}
	cookieState, err := cookies.ReadSigned(r, cookie_processor.CookieName, []byte(c.cfg.CookieSecret))
	if err != nil {
		return nil, nil, fmt.Errorf("cookies.ReadSigned: %w", err)
	}
	if cookieState != r.FormValue("state") {
		return nil, nil, fmt.Errorf("state is not generated by this Client: %w", err)
	}

	tokens, introspectResponse, err := c.oidc.Callback(r.Context(), req.Code)
	if err != nil {
		return nil, nil, fmt.Errorf("c.oidc.Callback: %w", err)
	}

	return tokens, introspectResponse, nil
}

type (
	CallbackRequest struct {
		Code  string
		State state.State
	}
)

func (r *CallbackRequest) Validate(req *http.Request) error {
	code := req.FormValue("code")
	st := req.FormValue("state")

	if code == "" {
		return errors.New("authorization code is empty")
	}
	if st == "" {
		return errors.New("state is empty")
	}

	stateModel, err := state.New(st)
	if err != nil {
		return err
	}

	r.Code = code
	r.State = stateModel

	return nil
}
