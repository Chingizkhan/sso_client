package oauth_oidc

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/sso_client/service/client"
	"log"
)

const (
	tokenIntrospectURL = "/oauth/token/introspect"
)

func Introspect(ctx context.Context, accessToken string) error {
	introspect, err := client.OauthClient.Introspect(ctx, tokenIntrospectURL, accessToken)
	if err != nil {
		return fmt.Errorf("OauthClient.Introspect: %w", err)
	}

	log.Println("oauth_oidc introspect response: ", introspect)

	return nil
}
