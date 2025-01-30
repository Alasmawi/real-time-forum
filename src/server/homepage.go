package server

import (
	"database/sql"
	fUtil "forum/database/feature-utils"
	"fmt"
	"log"
	"net/http"
	"time"
	strct "forum/structs"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("Method not allowed")
		err := strct.ErrorPageData{Code: "405", ErrorMsg: "METHOD NOT ALLOWED"}
		errHandler(w, r, &err)
		return
	}

	db, err := sql.Open("sqlite3", "./database/main.db")
	if err != nil {
		log.Println("Database connection failed:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}
	defer db.Close()

	var hasSession bool
	var userID int
	var userName string
	seshCok, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("No cookie found, treated as guest")
	} else if seshCok.Value == "" {
		hasSession = false
	} else {
		hasSession = true

		seshVal := seshCok.Value

		err = db.QueryRow("SELECT user_id, username FROM user WHERE current_session = ?", seshVal).Scan(&userID, &userName)
		if err == sql.ErrNoRows {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if err != nil {
			log.Println("Error fetching userid ID from user table:", err)
			err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &err)
			return
		} else {
			log.Println("User is logged in:", userName)
		}
	}

	users, err := fUtil.GetAllUsers(db)
	if err != nil {
		log.Println("Failed to fetch users:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	var posts []strct.Post
	posts, err = fUtil.GetAllPosts(db)
	if err != nil {
		log.Println("Failed to fetch posts:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	if hasSession {
		var avatar sql.NullString
		err = db.QueryRow("SELECT avatar FROM user WHERE user_id = ?", userID).Scan(&avatar)
		if err == sql.ErrNoRows {
			log.Println("No user found with the given ID:", userID)
			err := strct.ErrorPageData{Code: "404", ErrorMsg: "USER NOT FOUND"}
			errHandler(w, r, &err)
			return
		} else if err != nil {
			log.Println("Failed to fetch user data:", err)
			err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &err)
			return
		}

		var totalPosts int

		err = db.QueryRow("SELECT COUNT(*) FROM post WHERE user_id = ?", userID).Scan(&totalPosts)
		if err != nil {
			log.Println("Failed to fetch total posts:", err)
			err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &err)
			return
		}

		data := strct.PageData{
			HasSession:     hasSession,
			UserID:         userID,
			UserName:       userName,
			Avatar:         avatar.String,
			TotalPosts:     totalPosts,
			Users:          users,
			Posts:          posts,
			// Notifications:  notifications,
		}

		err = templates.ExecuteTemplate(w, "home.html", data)
		if err != nil {
			log.Println("Error rendering home page:", err)
			err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
			errHandler(w, r, &err)
			return
		}
	}
}
