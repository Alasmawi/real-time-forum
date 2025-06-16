package api

import (
	"net/http"
	"time"

	"reboot01.com/js/realtime-forum/internal/database"
	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/security"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete session from database

		app.DB.DeleteSession(cookie.Value)
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func (app *Application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Identifier string              `json:"identifier"`
		Password   string              `json:"password"`
		Validator  validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.Identifier != "", "Identifier", "Email or username is required")
	input.Validator.CheckField(input.Password != "", "Password", "Password is required")

	var user *database.User
	var found bool

	switch validator.IsEmail(input.Identifier) {
	case true:
		user, found, err = app.DB.GetUserByEmail(input.Identifier)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		input.Validator.CheckField(found, "Identifier", "Email address could not be found")
	case false:
		user, found, err = app.DB.GetUserByUsername(input.Identifier)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		input.Validator.CheckField(found, "Identifier", "Username could not be found")
	}

	if found && user != nil {
		passwordMatches, err := security.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	// Generate a new session token
	sessionToken, err := security.GenerateToken()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	stringToken := sessionToken.String()

	// Store session in database
	_, err = app.DB.InsertSession(user.ID, stringToken)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    stringToken,
		Expires:  time.Now().Add(30 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)

	// var claims jwt.Claims
	// claims.Subject = strconv.Itoa(user.ID)
	// expiry := time.Now().Add(24 * time.Hour)
	// claims.Issued = jwt.NewNumericTime(time.Now())
	// claims.NotBefore = jwt.NewNumericTime(time.Now())
	// claims.Expires = jwt.NewNumericTime(expiry)
	// claims.Issuer = app.config.baseURL
	// claims.Audiences = []string{app.config.baseURL}
	// jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// data := map[string]string{
	// "AuthenticationToken":       string(jwtBytes),
	// "AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	// }
	//write cookie here
	// err = response.JSON(w, http.StatusOK, data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}
