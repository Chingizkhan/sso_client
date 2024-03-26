package sso_service_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Chingizkhan/sso_client/pkg/token"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type (
	Client interface {
		Introspect(ctx context.Context, path, accessToken string) (*IntrospectResponse, error)
		Auth(ctx context.Context, path string, req *AuthRequest) (*AuthResponse, error)
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

func (s *OauthServiceClient) Auth(ctx context.Context, path string, in *AuthRequest) (*AuthResponse, error) {
	js, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal auth data: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.path+path, bytes.NewBuffer(js))
	if err != nil {
		return nil, fmt.Errorf("can not create http.Request: %w", err)
	}

	log.Println("s.path+path:", s.path+path)

	req.Header.Add("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()

	var response AuthResponse
	accessToken := resp.Header.Get("Access-Token")
	if accessToken == "" {
		return nil, fmt.Errorf("empty access token")
	}

	response.AccessToken = accessToken

	return &response, nil
}

func (s *OauthServiceClient) Introspect(ctx context.Context, path, accessToken string) (*IntrospectResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.path+path, nil)
	if err != nil {
		return nil, fmt.Errorf("can not create http.Request: %w", err)
	}

	log.Println("s.path+path:", s.path+path)

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
	log.Println("body:", string(body))
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	return nil
}
