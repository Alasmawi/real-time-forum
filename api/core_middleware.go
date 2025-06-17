package api

import (
	"fmt"
	"log/slog"
	"net/http"

	// "strconv"
	// "strings"
	// "time"

	"reboot01.com/js/realtime-forum/internal/response"
	// "github.com/pascaldekloe/jwt"
	// "github.com/tomasen/realip"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("response", "status", mw.StatusCode, "size", mw.BytesCount)

		app.Logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (app *Application) checkAuthenticated(w http.ResponseWriter, r *http.Request) {
	user := contextGetAuthenticatedUser(r)

	if user != nil {
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": true,
		})
		if err != nil {
			app.serverError(w, r, err)
		}
	} else {
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		if err != nil {
			app.serverError(w, r, err)
		}
	}
}

// authenticate retrieves the session token in the request cookie.
// If a valid session token is found, it retrieves the user associated with that session
// and adds it to the request context.
func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve session cookie and check if it exists
		cookie, err := r.Cookie("session_token")
		// If there is an error (cookie does not exist), err == nil is false
		// If there is no error (cookie exists), err == nil is true
		hasCookie := err == nil

		switch hasCookie {
		case true:
			if cookie.Value == "" {
				app.invalidateSessionToken(w, r)
				return
			}

			user, found, err := app.DB.GetUserBySession(cookie.Value)
			if err != nil {
				app.serverError(w, r, err)
				return
			}

			if !found {
				app.invalidateSessionToken(w, r)
				return
			}

			// Adds users with valid sessions to the context
			r = contextSetAuthenticatedUser(r, user)

		case false:
			// Potential guest logic can be added here
		}

		next.ServeHTTP(w, r)
	})
}

// Authorize requires an authenticated user in context to grant access to the next handler.
// If the user is not authenticated, it responds with a 401 Unauthorized status.
func (app *Application) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := contextGetAuthenticatedUser(r)

		if authenticatedUser == nil {
			app.authenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//func (app *Application) authenticate(next http.Handler) http.Handler {
// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Vary", "Authorization")

// 	authorizationHeader := r.Header.Get("Authorization")

// 	if authorizationHeader != "" {
// 		headerParts := strings.Split(authorizationHeader, " ")

// 		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
// 			token := headerParts[1]

// 			claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secretKey))
// 			if err != nil {
// 				app.invalidAuthenticationToken(w, r)
// 				return
// 			}

// 			if !claims.Valid(time.Now()) {
// 				app.invalidAuthenticationToken(w, r)
// 				return
// 			}

// 			if claims.Issuer != app.config.baseURL {
// 				app.invalidAuthenticationToken(w, r)
// 				return
// 			}

// 			if !claims.AcceptAudience(app.config.baseURL) {
// 				app.invalidAuthenticationToken(w, r)
// 				return
// 			}

// 			userID, err := strconv.Atoi(claims.Subject)
// 			if err != nil {
// 				app.serverError(w, r, err)
// 				return
// 			}

// 			user, found, err := app.db.GetUser(userID)

// 			if err != nil {
// 				app.serverError(w, r, err)
// 				return
// 			}

// 			if found {
// 				r = contextSetAuthenticatedUser(r, user)
// 			}
// 		}
// 	}

// 	next.ServeHTTP(w, r)
// })
//}

// func (app *Application) requireAuthenticatedUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authenticatedUser := contextGetAuthenticatedUser(r)

// 		if authenticatedUser == nil {
// 			app.authenticationRequired(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
