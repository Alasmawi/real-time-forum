package cookie

import (
	"net/http"
	"time"
)

const CookieExpirey = 30 * time.Minute

func SetSessionCookie(w http.ResponseWriter, name, value, path, domain string, httpOnly, secure bool, sameSite http.SameSite, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   maxAge,
	})
}

func SetDefaultSessionCookie(w http.ResponseWriter, value string) {
	SetSessionCookie(w, "session_token", value, "/", "", false, false, http.SameSiteStrictMode, int((CookieExpirey).Seconds()))
}

func ClearDefaultSessionCookie(w http.ResponseWriter) {
	SetSessionCookie(w, "session_token", "", "/", "", false, false, http.SameSiteStrictMode, -1)
}
