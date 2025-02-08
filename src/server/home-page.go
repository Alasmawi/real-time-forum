package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	fUtil "reboot01.com/js/forum/database/feature-utils"
	strct "reboot01.com/js/forum/structs"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	var hasSession bool
	var userID int
	var userName string

	seshCok, _ := r.Cookie("session_token")

	seshVal := seshCok.Value

	err := strct.Db.QueryRow("SELECT user_id, username FROM user WHERE current_session = ?", seshVal).Scan(&userID, &userName)
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
		err := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	} else {
		log.Println("User is logged in:", userName)
	}

	users, err := fUtil.GetAllUsers()
	if err != nil {
		log.Println("Failed to fetch users:", err)
		err := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	var posts []strct.Post
	posts, err = fUtil.GetAllPosts()
	if err != nil {
		log.Println("Failed to fetch posts:", err)
		err := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	data := strct.PageData{
		HasSession: hasSession,
		UserID:     userID,
		UserName:   userName,
		Users:      users,
		Posts:      posts,
		// Notifications:  notifications,
	}

	err = strct.Templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		log.Println("Error rendering home page:", err)
		err := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}
}
