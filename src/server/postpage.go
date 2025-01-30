package server

import (
	"database/sql"
	"fmt"
	fUtils "forum/database/feature-utils"
	strct "forum/structs"
	"log"
	"net/http"
	"strconv"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		log.Println("Invalid URL path")
		err := strct.ErrorPageData{Code: "404", ErrorMsg: "PAGE NOT FOUND"}
		errHandler(w, r, &err)
		return
	}

	if r.Method != "GET" {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./database/main.db")
	if err != nil {
		log.Println("Error opening database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Fetch session cookie
	seshCok, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		fmt.Println("Error fetching session cookie")
		return
	}

	// Set session token from cookie value
	seshVal := seshCok.Value

	var userID int
	var userName string
	err = db.QueryRow("SELECT user_id, username FROM user WHERE current_session = ?", seshVal).Scan(&userID, &userName)
	if err != nil {
		log.Println("Error fetching userid and username from user table:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		log.Println("Post ID not found in query parameters")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	var post strct.Post
	err = db.QueryRow(`
        SELECT post.post_id, post.image, post.content, post.post_at, post.user_id, user.username, user.f_name, user.l_name, user.avatar,
        (SELECT COUNT(*) FROM comment WHERE comment.post_id = post.post_id) AS Comments
        FROM post
        JOIN user ON post.user_id = user.user_id
        WHERE post.post_id = ?
		`, postID).Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Comments)
	if err != nil {
		log.Println("Failed to fetch posts")
		errData := strct.ErrorPageData{Code: "400", ErrorMsg: "BAD REQUEST"}
		errHandler(w, r, &errData)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		log.Println("Error converting post ID to integer:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	comments, err := fUtils.GetCommentsForPost(db, postIDInt)
	if err != nil {
		log.Println("Error getting comments for post:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if userID < 0 {
		log.Println("User ID not found in query parameters")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// Fetch categories for the post
	categories, err := fUtils.GetCategoriesForPost(db, post.PostID)
	if err != nil {
		log.Println("Error fetching categories for post:", err)
		return
	}

	log.Println("User ID:", userID) // Log the UserID to ensure it is being retrieved

	data := strct.PageData{
		Post:       post,
		Comments:   comments,
		UserID:     userID, // Ensure UserID is set
		UserName:   userName,
		Categories: categories,
	}

	err = templates.ExecuteTemplate(w, "post.html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
