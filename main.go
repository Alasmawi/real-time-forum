package main

import (
	"fmt"
	db "reboot01.com/js/forum/database/sql-utils"
	"log"
	"net/http"

	s "reboot01.com/js/forum/src/server"
)

func init() {
	// db.DropDataBase()
	db.DataBase()
}

const port = ":8080"

func main() {
	//make fetching userID with global variable
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", s.ReverseMiddleware(s.LoginPage))
	mux.HandleFunc("POST /logout", s.AuthMiddleware(s.Logout))
	mux.HandleFunc("/signup", s.SignupPage)
	mux.HandleFunc("GET /home", s.AuthMiddleware(s.HomePage))
	mux.HandleFunc("/create", s.AuthMiddleware(s.CreatePage))
	mux.HandleFunc("GET /notifications", s.AuthMiddleware(s.NotificationsPage))
	mux.HandleFunc("GET /post/{postID}", s.AuthMiddleware(s.PostPage))
	mux.HandleFunc("POST /deletepost", s.AuthMiddleware(s.DeletePost))
	mux.HandleFunc("POST /addcomment", s.AuthMiddleware(s.AddComment))

	fmt.Println("Server running on http://localhost:8080/\nTo stop the server press Ctrl+C")

	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
