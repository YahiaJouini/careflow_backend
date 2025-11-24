package auth

import "net/http"

type CookieMode string

const (
	Remove CookieMode = "remove"
	Add    CookieMode = "add"
)

type MobileAuthResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         interface{} `json:"user"`
}

func SetAuthCookie(w http.ResponseWriter, value string, mode CookieMode) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Path:     "/",
		Value:    value,
		HttpOnly: true,
		Secure:   false, // Change to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	if mode == Remove {
		cookie.MaxAge = -1 // Expire immediately
	} else {
		cookie.MaxAge = 60 * 60 * 24 * 30 // 30 days
	}

	http.SetCookie(w, cookie)
}

func Logout(w http.ResponseWriter) {
	SetAuthCookie(w, "", Remove)
}
