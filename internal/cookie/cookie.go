package cookie

import (
	"net/http"
	"time"
)

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
	SetSessionCookie(w, "session_token", value, "/", "", true, false, http.SameSiteStrictMode, int((30 * time.Minute).Seconds()))
}

func ClearDefaultSessionCookie(w http.ResponseWriter) {
	SetSessionCookie(w, "session_token", "", "/", "", true, false, http.SameSiteStrictMode, -1)
}
