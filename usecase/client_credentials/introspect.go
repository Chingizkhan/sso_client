package client_credentials

import (
	"context"
	"fmt"
	"log"
	"sso_client/service/client"
)

const (
	tokenIntrospectURL = "/client-credentials/token/introspect"
)

func Introspect(ctx context.Context, accessToken string) error {
	introspect, err := client.OauthClient.Introspect(ctx, tokenIntrospectURL, accessToken)
	if err != nil {
		return fmt.Errorf("OauthClient.Introspect: %w", err)
	}

	log.Println("client_credentials introspect response: ", introspect)

	return nil
}
