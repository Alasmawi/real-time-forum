package server

import (
	strct "forum/structs"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("Method not allowed")
		err := strct.ErrorPageData{Code: "405", ErrorMsg: "METHOD NOT ALLOWED"}
		errHandler(w, r, &err)
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		log.Println("Post ID is missing")
		err := strct.ErrorPageData{Code: "400", ErrorMsg: "BAD REQUEST"}
		errHandler(w, r, &err)
		return
	}
	_, err := strct.Db.Exec("DELETE FROM post WHERE post_id = ?", postID)
	if err != nil {
		log.Println("Error deleting post:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}