package api

import (
	"net/http"
	"time"

	// "strconv"
	// "time"

	"reboot01.com/js/realtime-forum/internal/database"
	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/response"
	"reboot01.com/js/realtime-forum/internal/security"
	"reboot01.com/js/realtime-forum/internal/validator"
	// "github.com/pascaldekloe/jwt"
)

// func (app *Application) serverIndex() http.Handler {

// }

func (app *Application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) serverIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string              `json:"username"`
		Email     string              `json:"email"`
		Password  string              `json:"password"`
		Age       string              `json:"age"`
		Sex       string              `json:"sex"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, usernameFound, err := app.DB.GetUserByUsername(input.Username)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, emailFound, err := app.DB.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Username != "", "Username", "Username is required")
	input.Validator.CheckField(!usernameFound, "Username", "Username is already in use")

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	input.Validator.CheckField(!emailFound, "Email", "Email is already in use")

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, security.CommonPasswords...), "Password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := security.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.DB.InsertUser(input.Email, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		Expires:  time.Now().Add(30 * time.Second),
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

// func (app *Application) protected(w http.ResponseWriter, r *http.Request) {
// w.Write([]byte("This is a protected handler"))
// }

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

func (app *Application) checkAuthentication(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No session cookie
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"redirect":      "/",
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
			"redirect":      "/",
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	// Validate session exists in database and get user
	user, found, err := app.DB.GetUserBySession(cookie.Value)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If session doesn't exist in database
	if !found {
		// Delete invalid session from database if it exists
		app.DB.DeleteSession(cookie.Value)
		err := response.JSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"redirect":      "/",
		})
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	// Valid session - return user info
	err = response.JSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
		"redirect": "/posts",
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}
