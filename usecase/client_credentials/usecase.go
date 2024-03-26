package client_credentials

import (
	"context"
	"fmt"
	"github.com/Chingizkhan/sso_client/service/sso_service_client"
)

const (
	tokenIntrospectURL = "/client-credentials/token/introspect"
	tokenAuthURL       = "/client-credentials/auth"
)

type (
	UseCase struct {
		client sso_service_client.Client
	}
)

func New(client sso_service_client.Client) *UseCase {
	return &UseCase{client}
}

func (u *UseCase) Introspect(ctx context.Context, accessToken string) (*sso_service_client.IntrospectResponse, error) {
	introspect, err := u.client.Introspect(ctx, tokenIntrospectURL, accessToken)
	if err != nil {
		return nil, fmt.Errorf("OauthClient.Introspect: %w", err)
	}

	return introspect, nil
}

func (u *UseCase) Auth(ctx context.Context, in *sso_service_client.AuthRequest) (*sso_service_client.AuthResponse, error) {
	introspect, err := u.client.Auth(ctx, tokenAuthURL, in)
	if err != nil {
		return nil, fmt.Errorf("OauthClient.Auth: %w", err)
	}

	return introspect, nil
}
