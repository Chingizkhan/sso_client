package oauth_oidc

//var (
//	CookieName = "oauth_state"
//)
//
//func LoginFlow(w http.ResponseWriter, secret string) (string, error) {
//	st, err := state.Generate()
//	if err != nil {
//		return "", fmt.Errorf("state.Generate: %w", err)
//	}
//
//	// todo: think about cookie service (set secret, expiration_time, methods, etc...)
//	expiration := time.Now().Add(10 * time.Minute)
//	cookie := http.Cookie{
//		Name:     CookieName,
//		Value:    string(st),
//		Expires:  expiration,
//		Secure:   true,
//		Path:     "/",
//		HttpOnly: true,
//		SameSite: http.SameSiteLaxMode,
//	}
//	err = cookies.WriteSigned(w, cookie, []byte(secret))
//	if err != nil {
//		return "", fmt.Errorf("cookies.WriteSigned: %w", err)
//	}
//
//	loginUrl := oauth.Config.AuthCodeURL(string(st))
//
//	return loginUrl, nil
//}
