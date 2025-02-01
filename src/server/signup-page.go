package server

import (
	strct "forum/structs"
	"log"
	"net/http"
	"regexp"
)

func SignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		F_name := r.FormValue("first_name")
		L_name := r.FormValue("last_name")
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
		if !emailRegex.MatchString(email) {
			err := strct.Templates.ExecuteTemplate(w, "signup.html", map[string]string{
				"ErrorMessage": "Invalid email format",
			})
			if err != nil {
				log.Println("Error rendering signup page:", err)
				errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
			}
			return
		}

		if password != confirmPassword {
			err := strct.Templates.ExecuteTemplate(w, "signup.html", map[string]string{
				"ErrorMessage": "Passwords do not match",
			})
			if err != nil {
				log.Println("Error rendering signup page:", err)
				errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
			}
			return
		}

		var usernameExists, emailExists bool
		err := strct.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", username).Scan(&usernameExists)
		if err != nil {
			log.Println("Failed to check if username exists")
			errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		}

		err = strct.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", email).Scan(&emailExists)
		if err != nil {
			log.Println("Failed to check if email exists")
			errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		}

		if usernameExists {
			err := strct.Templates.ExecuteTemplate(w, "signup.html", map[string]string{
				"ErrorMessage": "Username already exists",
			})
			if err != nil {
				log.Println("Error rendering signup page:", err)
				errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
			}
			return
		}

		if emailExists {
			err := strct.Templates.ExecuteTemplate(w, "signup.html", map[string]string{
				"ErrorMessage": "Email already exists",
			})
			if err != nil {
				log.Println("Error rendering signup page:", err)
				errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
			}
			return
		}
		password, _ = HashPassword(password)
		defaultAvatar := "static/assets/default-avatar.png"
		// Insert user data into the database
		stmt, err := strct.Db.Prepare("INSERT INTO user (f_name, l_name, username, email, password, current_session, avatar) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Println("Failed to prepare insert statement:", err)
			errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(F_name, L_name, username, email, password, "", defaultAvatar)
		if err != nil {
			log.Println("Failed to insert user data:", err)
			errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		}

		// Redirect to login page or show success message
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := strct.Templates.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		log.Println("Error rendering signup page:", err)
		errData := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &errData)
	}
}
