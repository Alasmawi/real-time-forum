package utilities

import (
	"fmt"
	"net/http"
)

func FetchCookieValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}

	return cookie.Value, nil
}
