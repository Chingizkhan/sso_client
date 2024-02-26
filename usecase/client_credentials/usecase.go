package client_credentials

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/sso_client/service/sso_service_client"
	"log"
)

type (
	UseCase struct {
		client sso_service_client.Client
	}
)

func New(client sso_service_client.Client) *UseCase {
	return &UseCase{client}
}

func (u *UseCase) Introspect(ctx context.Context, accessToken string) error {
	introspect, err := u.client.Introspect(ctx, tokenIntrospectURL, accessToken)
	if err != nil {
		return fmt.Errorf("OauthClient.Introspect: %w", err)
	}

	log.Println("client_credentials introspect response: ", introspect)

	return nil
}
