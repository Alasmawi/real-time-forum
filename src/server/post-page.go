package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	fUtils "reboot01.com/js/forum/database/feature-utils"
	strct "reboot01.com/js/forum/structs"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
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
	err = strct.Db.QueryRow("SELECT user_id, username FROM user WHERE current_session = ?", seshVal).Scan(&userID, &userName)
	if err != nil {
		log.Println("Error fetching userid and username from user table:", err)
		err := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	postID := r.PathValue("postID")
	if postID == "" {
		log.Println("Post ID not found in query parameters")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	var post strct.Post
	err = strct.Db.QueryRow(`
        SELECT post.post_id, post.image, post.content, post.post_at, post.user_id, user.username, user.f_name, user.l_name, user.avatar,
        (SELECT COUNT(*) FROM comment WHERE comment.post_id = post.post_id) AS Comments
        FROM post
        JOIN user ON post.user_id = user.user_id
        WHERE post.post_id = ?
		`, postID).Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Comments)
	if err != nil {
		log.Println("Failed to fetch posts", err)
		errData := strct.ErrorPageData{Code: http.StatusInternalServerError, ErrorMsg: "BAD REQUEST"}
		errHandler(w, r, &errData)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		log.Println("Error converting post ID to integer:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	comments, err := fUtils.GetCommentsForPost(postIDInt)
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
	categories, err := fUtils.GetCategoriesForPost(post.PostID)
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

	err = strct.Templates.ExecuteTemplate(w, "post.html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
