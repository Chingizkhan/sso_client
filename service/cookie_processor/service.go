package cookie_processor

import (
	"net/http"
	"time"
)

var (
	CookieName = "oauth_state"
)

type CookieProcessor struct {
	cookieLifetime time.Duration
}

func New(lifetime time.Duration) *CookieProcessor {
	return &CookieProcessor{
		cookieLifetime: lifetime,
	}
}

func (p *CookieProcessor) GenerateCookie(state string) http.Cookie {
	expiration := time.Now().Add(p.cookieLifetime)
	cookie := http.Cookie{
		Name:     CookieName,
		Value:    state,
		Expires:  expiration,
		Secure:   true,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}
