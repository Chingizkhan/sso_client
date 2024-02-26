package sso_client

import (
	"github.com/Chingizkhan/sso_client/pkg/api_util"
	"github.com/Chingizkhan/sso_client/pkg/cookies"
	"github.com/Chingizkhan/sso_client/pkg/token"
	"github.com/Chingizkhan/sso_client/service/cookie_processor"
	"github.com/Chingizkhan/sso_client/service/sso_service_client"
	"github.com/Chingizkhan/sso_client/usecase/client_credentials"
	"github.com/Chingizkhan/sso_client/usecase/oauth_oidc"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

type (
	Client interface {
		AuthOauth2(next http.Handler) http.Handler
		AuthClientCredentials(next http.Handler) http.Handler
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
	//oauth2Config := oauth2.Config{
	//	ClientID:     "38b36b9d-48a8-40fd-9911-ee4462428c58",
	//	ClientSecret: "mysecret",
	//	RedirectURL:  "http://localhost:8082/callback",
	//	Scopes:       []string{"offline", "users.write", "users.read", "users.edit", "users.delete"},
	//	Endpoint: oauth2.Endpoint{
	//		AuthURL:  "http://localhost:9010/oauth2/auth",
	//		TokenURL: "http://localhost:9010/oauth2/token",
	//	},
	//}
	cookieProcessor := cookie_processor.New(config.CookieLifetime)
	ssoServiceClient := sso_service_client.New(time.Second*15, config.OauthAddr)
	oidc := oauth_oidc.New(config.OauthAddr, cookieProcessor, ssoServiceClient, config.Oauth2Config)
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
		log.Println("accessToken", accessToken)
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
		if err := c.oidc.Introspect(r.Context(), accessToken); err != nil {
			api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
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
		if err := c.clientCredentials.Introspect(r.Context(), accessToken); err != nil {
			api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
