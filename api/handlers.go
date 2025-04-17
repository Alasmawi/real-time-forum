package api

import (
	"net/http"
	// "strconv"
	// "time"

	"reboot01.com/js/forum/internal/database"
	"reboot01.com/js/forum/internal/password"
	"reboot01.com/js/forum/internal/request"
	"reboot01.com/js/forum/internal/response"
	"reboot01.com/js/forum/internal/validator"
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

func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username   string              `json:"Username"`
		Identifier string              `json:"Identifier"`
		Password   string              `json:"Password"`
		Age        string              `json:"Age"`
		Sex        string              `json:"Sex"`
		Validator  validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, found, err := app.DB.GetUserByEmail(input.Identifier)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Identifier != "", "Identifier", "Identifier is required")
	input.Validator.CheckField(validator.Matches(input.Identifier, validator.RgxEmail), "Identifier", "Must be a valid email address")
	input.Validator.CheckField(!found, "Identifier", "Identifier is already in use")

	input.Validator.CheckField(input.Identifier != "", "Identifier", "Identifier is required")
	input.Validator.CheckField(validator.Matches(input.Identifier, validator.RgxEmail), "Identifier", "Must be a valid email address")
	input.Validator.CheckField(!found, "Identifier", "Identifier is already in use")

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.DB.InsertUser(input.Identifier, hashedPassword)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Identifier string              `json:"Identifier"`
		Password   string              `json:"Password"`
		Validator  validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var user *database.User
	var found bool

	input.Validator.CheckField(input.Identifier != "", "Identifier", "Identifier is required")
	input.Validator.CheckField(found, "Identifier", "Identifier address could not be found")

	if validator.IsEmail(input.Identifier) {
		user, found, err = app.DB.GetUserByEmail(input.Identifier)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		user, found, err = app.DB.GetUserByUsername(input.Identifier)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

	}

	if found {
		passwordMatches, err := password.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	

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

	data := map[string]string{
		// "AuthenticationToken":       string(jwtBytes),
		// "AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	//write cookie here

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
