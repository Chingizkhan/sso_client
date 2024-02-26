package middleware

import (
	"github.com/Chingizkhan/sso_client/pkg/api_util"
	"github.com/Chingizkhan/sso_client/pkg/token"
	"github.com/Chingizkhan/sso_client/usecase/oauth_oidc"
	"log"
	"net/http"
)

func AuthOauth2(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// check token exists
			accessToken, err := token.FromHttpRequest(r)
			if err != nil {
				api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			log.Println("accessToken", accessToken)
			if accessToken == "" {
				// return login page
				loginURL, err := oauth_oidc.LoginFlow(w, r, secret)
				if err != nil {
					api_util.RenderErrorResponse(w, "can not show login flow", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
				// todo: redirect here. not in LoginFLow
				return
			}

			// introspect token
			if err := oauth_oidc.Introspect(r.Context(), accessToken); err != nil {
				api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
