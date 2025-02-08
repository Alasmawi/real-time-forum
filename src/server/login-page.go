package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	uuid "reboot01.com/js/forum/src/uuid"
	strct "reboot01.com/js/forum/structs"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := strct.Templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Println("Error rendering login page:", err)
			errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
		}

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		var userID int
		var dbPassword, userName string
		err := strct.Db.QueryRow("SELECT user_id, password, username FROM user WHERE email = ?", email).Scan(&userID, &dbPassword, &userName)
		if err != nil {
			if err == sql.ErrNoRows {
				// No credentials found with the given email
				err = strct.Templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
					"ErrorMsg": "Invalid email or password",
				})
				if err != nil {
					log.Println("Error rendering login page:", err)
					errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
					errHandler(w, r, &errData)
				}
				return
			}
			log.Println("Failed to fetch user data")
			err := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &err)
			return
		}

		if !VerifyPassword(password, dbPassword) {
			err := strct.Templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
				"ErrorMsg": "Invalid email or password",
			})
			if err != nil {
				log.Println("Error rendering login page:", err)
				errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
			}
			return
		}

		// Generate a new session token
		sessionToken, err := uuid.GenerateToken()
		if err != nil {
			log.Println("Error generating UUID:", err)
			errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
		}

		stringToken := sessionToken.String()

		//Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    stringToken,
			Expires:  time.Now().Add(1 * time.Hour), //1 hour lifetime
			HttpOnly: true,
		})

		// Update the user's session ID in the session table
		result, err := strct.Db.Exec("UPDATE session SET session_id = ? WHERE user_id = ?", stringToken, userID)
		if err != nil {
			log.Println("Error updating session ID in session table:", err)
			errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		} else if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected == 0 { //only insert a new row if no record is updated (i.e., no session is found)
			_, err := strct.Db.Exec("INSERT INTO session (session_id, user_id, end_time) VALUES (?, ?, ?) RETURNING session_id",
				stringToken, userID, time.Now().Add(1*time.Hour))
			if err != nil {
				log.Println("Error creating new session:", err)
				errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &errData)
				return
			}
		}

		// Update the user's session ID in the database
		_, err = strct.Db.Exec("UPDATE user SET current_session = ? WHERE user_id = ?", stringToken, userID)
		if err != nil {
			log.Println("Error updating session ID in user table:", err)
			errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &errData)
			return
		}

		log.Println("User logged in with userID:", userID)

		log.Println("Redirecting to home page")
		http.Redirect(w, r, "/home", http.StatusSeeOther)

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		errData := strct.ErrorPageData{Code: http.StatusMethodNotAllowed, ErrorMsg: "METHOD NOT ALLOWED"}
		errHandler(w, r, &errData)
	}
}
