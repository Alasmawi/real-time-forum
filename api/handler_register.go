package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/database"
	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/security"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
		F_name    string              `json:"f_name"`
		L_name    string              `json:"l_name"`
		Username  string              `json:"username"`
		Email     string              `json:"email"`
		Password  string              `json:"password"`
		Age       int                 `json:"age"`
		Sex       bool                `json:"sex"`
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

	input.Validator.CheckField(input.F_name != "", "First-name", "First name is required")
	input.Validator.CheckField(validator.MaxRunes(input.F_name, 15), "First-name", "First name must be at most 15 characters long")
	input.Validator.CheckField(input.L_name != "", "Last-name", "Last name is required")
	input.Validator.CheckField(validator.MaxRunes(input.L_name, 15), "Last-name", "Last name must be at most 15 characters long")

	input.Validator.CheckField(input.Age > 0 && input.Age <= 120, "Age", "Age must be a realistic value (0-120)")

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

	user := database.User{
		FName:          input.F_name,
		LName:          input.L_name,
		Username:       input.Username,
		Email:          input.Email,
		HashedPassword: hashedPassword,
		Sex:            input.Sex,
		Age:            input.Age,
	}

	_, err = app.DB.InsertUser(&user)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
