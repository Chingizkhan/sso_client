package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"sso_client/config"
	"sso_client/pkg/token"
	"strings"
	"time"
)

var OauthClient *OauthServiceClient

func init() {
	OauthClient = New(time.Second*15, config.Cfg.OauthService.Addr)
}

type (
	Client interface {
		Introspect(ctx context.Context, path, accessToken string) (IntrospectResponse, error)
	}

	OauthServiceClient struct {
		client *http.Client
		path   string
	}
)

func New(timeout time.Duration, path string) *OauthServiceClient {
	return &OauthServiceClient{
		client: &http.Client{Timeout: timeout},
		path:   path,
	}
}

func (s *OauthServiceClient) Introspect(ctx context.Context, path, accessToken string) (*IntrospectResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.path+path, nil)
	if err != nil {
		return nil, fmt.Errorf("can not create http.Request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.TypeBearer, accessToken))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if strings.Contains(string(body), "error") {
		var res ErrorResponse
		if err = getResponse(body, &res); err != nil {
			return nil, fmt.Errorf("get 'error' response: %w", err)
		}
		return nil, errors.New(res.Error)
	}

	var response IntrospectResponse
	if err = getResponse(body, &response); err != nil {
		return nil, fmt.Errorf("get response: %w", err)
	}

	return &IntrospectResponse{
		Active:   response.Active,
		ClientID: response.ClientID,
		Exp:      response.Exp,
		Sub:      response.Sub,
		UserName: response.UserName,
		TokenUse: response.TokenUse,
	}, nil
}

func getResponse(body []byte, res any) error {
	err := json.Unmarshal(body, res)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	return nil
}