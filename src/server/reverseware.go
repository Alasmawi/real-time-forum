package server

import (
	"fmt"
	strct "forum/structs"
	"log"
	"net/http"
	"time"
)

func ReverseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fetch session cookie
		seshCok, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("This user has no cookie")
		} else {
			//Set session token from cookie value
			seshVal := seshCok.Value

			var exists bool
			err = strct.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE current_session = ?)", seshVal).Scan(&exists)
			if err != nil {
				log.Println("Error :", err)
				err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
				errHandler(w, r, &err)
				return
			}

			if !exists {
				log.Println("Inavlid Session")
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    "",
					Expires:  time.Now().Add(-time.Hour),
					HttpOnly: true,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else if exists {
				fmt.Println("Valid session")
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}
		}
		next.ServeHTTP(w, r)
	})
}
