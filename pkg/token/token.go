package token

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrInvalidAuthHeader   = errors.New("invalid authorization header")
	ErrUnsupportedAuthType = errors.New("unsupported authorization type: must be 'Bearer'")
)

const (
	TypeBearer       = "Bearer"
	authorizationKey = "Authorization"
)

func FromHttpRequest(r *http.Request) string {
	return strings.TrimSpace(strings.Replace(r.Header.Get(authorizationKey), TypeBearer, "", 1))
}

//func FromHttpRequest(r *http.Request) (string, error) {
//	authKey := r.Header.Get(authorizationKey)
//	fields := strings.Fields(authKey)
//	if len(fields) < 2 {
//		return "", ErrInvalidAuthHeader
//	}
//	authType := strings.ToLower(fields[0])
//	if authType != TypeBearer {
//		return "", ErrUnsupportedAuthType
//	}
//	return fields[1], nil
//}
