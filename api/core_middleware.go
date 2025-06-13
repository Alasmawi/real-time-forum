package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

// authentication middleware
func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No session cookie, redirect to login
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Check if session value is empty
		if cookie.Value == "" {
			// Invalid session, clear cookie and redirect
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Validate session exists in database and get user
		user, found, err := app.DB.GetUserBySession(cookie.Value)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// If session doesn't exist in database, clear cookie and redirect
		if !found {
			// Delete invalid session from database if it exists
			app.DB.DeleteSession(cookie.Value)
			// Clear cookie and redirect
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Valid session: set user in context and continue
		r = contextSetAuthenticatedUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) checkAuthentication(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No session cookie
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	// Check if session value is empty
	if cookie.Value == "" {
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	// Validate session exists in database and get user
	_, found, err := app.DB.GetUserBySession(cookie.Value)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If session doesn't exist in database
	if !found {
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	// Valid session - return user info
	err = response.JSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}

// tokenAuthenticate middleware for API endpoints - uses session cookies with JSON error responses
func (app *Application) tokenAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No session cookie
			app.invalidateSessionToken(w, r)
			return
		}

		// Check if session value is empty
		if cookie.Value == "" {
			// Invalid session, clear cookie and return error
			app.invalidateSessionToken(w, r)
			return
		}

		// Validate session exists in database and get user
		user, found, err := app.DB.GetUserBySession(cookie.Value)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// If session doesn't exist in database, clean up and return error
		if !found {
			// Delete invalid session from database if it exists
			app.DB.DeleteSession(cookie.Value)
			app.invalidateSessionToken(w, r)
			return
		}

		// Valid session: set user in context and continue
		r = contextSetAuthenticatedUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// requireSessionUser middleware - requires a valid authenticated user in context for API endpoints
func (app *Application) requireSessionUser(next http.Handler) http.Handler {
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
