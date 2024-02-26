package middleware

import (
	"github.com/Chingizkhan/sso_client/pkg/api_util"
	"github.com/Chingizkhan/sso_client/pkg/token"
	"github.com/Chingizkhan/sso_client/usecase/client_credentials"
	"net/http"
)

func AuthClientCredentials() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// check token exists
			tok := token.FromHttpRequest(r)
			if tok == "" {
				api_util.RenderErrorResponse(w, "empty token", http.StatusBadRequest)
				return
			}

			// introspect token
			if err := client_credentials.Introspect(r.Context(), tok); err != nil {
				api_util.RenderErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
