package server

import (
	strct "forum/structs"
	"log"
	"net/http"
	"time"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("post_id")
	userID := r.FormValue("user_id")
	content := r.FormValue("content")

	log.Println("post_id:", postID)
	log.Println("user_id:", userID)
	log.Println("content:", content)

	if postID == "" || userID == "" || content == "" {
		log.Println("Missing form values")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err := strct.Db.Exec("INSERT INTO comment (content, comment_at, post_id, user_id) VALUES (?, ?, ?, ?)", content, time.Now(), postID, userID)
	if err != nil {
		log.Println("Error inserting comment:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
