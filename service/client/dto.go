package client

type (
	IntrospectResponse struct {
		Active   bool   `json:"active"`
		ClientID string `json:"client_id"`
		Exp      int64  `json:"exp"`
		Sub      string `json:"sub"`
		UserName string `json:"user_name"`
		TokenUse string `json:"token_use"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)
