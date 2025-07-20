package api

import (
	"net/http"
	"strconv"

	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/security"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	// Check if user is already authenticated
	if user := contextGetAuthenticatedUser(r); user != nil {
		app.Conflict(w, r)
		return
	}

	var input struct {
		FirstName string              `json:"first-name"`
		LastName  string              `json:"last-name"`
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

	// First Name validation
	input.Validator.CheckField(validator.NotBlank(input.FirstName), "first-name", "First name is required")
	input.Validator.CheckField(validator.MinRunes(input.FirstName, 2), "first-name", "First name must be at least 2 characters")
	input.Validator.CheckField(validator.MaxRunes(input.FirstName, 50), "first-name", "First name limit exceeded (50 characters)")

	// Last Name validation
	input.Validator.CheckField(validator.NotBlank(input.LastName), "last-name", "Last name is required")
	input.Validator.CheckField(validator.MinRunes(input.LastName, 2), "last-name", "Last name must be at least 2 characters")
	input.Validator.CheckField(validator.MaxRunes(input.LastName, 50), "last-name", "Last name limit exceeded (50 characters)")

	// Username validation
	input.Validator.CheckField(validator.NotBlank(input.Username), "username", "Username is required")
	input.Validator.CheckField(validator.MinRunes(input.Username, 3), "username", "Username must be at least 3 characters")
	input.Validator.CheckField(validator.MaxRunes(input.Username, 30), "username", "Username limit exceeded (30 characters)")
	input.Validator.CheckField(!usernameFound, "username", "Username is already in use")

	// Email validation
	input.Validator.CheckField(validator.NotBlank(input.Email), "email", "Email is required")
	input.Validator.CheckField(validator.IsEmail(input.Email), "email", "Must be a valid email address")
	input.Validator.CheckField(!emailFound, "email", "Email is already in use")

	// Password validation
	input.Validator.CheckField(validator.NotBlank(input.Password), "password", "Password is required")
	input.Validator.CheckField(validator.MinRunes(input.Password, 8), "password", "Password must be at least 8 characters")
	input.Validator.CheckField(validator.MaxRunes(input.Password, 72), "password", "Password must be no more than 72 characters")
	input.Validator.CheckField(validator.NotIn(input.Password, security.CommonPasswords...), "password", "Password is too common")

	// Age validation
	input.Validator.CheckField(validator.NotBlank(input.Age), "age", "Age is required")

	age, err := strconv.Atoi(input.Age)
	if err != nil {
		input.Validator.AddFieldError("age", "Age must be a valid number")
	}
	input.Validator.CheckField(validator.MinInt(age, 13), "age", "You must be 13 or older to register")
	input.Validator.CheckField(validator.MaxInt(age, 120), "age", "Dead men tell no tales...")

	// Sex validation
	input.Validator.CheckField(validator.NotBlank(input.Sex), "sex", "Sex is required")

	sex, err := strconv.Atoi(input.Sex)
	if err != nil {
		input.Validator.AddFieldError("sex", "Please select a valid option")
	}
	input.Validator.CheckField(validator.Between(sex, 0, 1), "sex", "Please select a valid option")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := security.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.DB.InsertUser(input.FirstName, input.LastName, input.Username, input.Email, hashedPassword, age, sex == 1)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
