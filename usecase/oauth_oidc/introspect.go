package oauth_oidc

import (
	"context"
	"fmt"
	"log"
	"sso_client/service/client"
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
