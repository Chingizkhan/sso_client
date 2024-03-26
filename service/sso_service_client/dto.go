package sso_service_client

type (
	IntrospectResponse struct {
		Active   bool   `json:"active"`
		ClientID string `json:"client_id"`
		Exp      int64  `json:"exp"`
		Sub      string `json:"sub"`
		UserName string `json:"user_name"`
		TokenUse string `json:"token_use"`
	}

	AuthRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	AuthResponse struct {
		AccessToken string `json:"access_token"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)
