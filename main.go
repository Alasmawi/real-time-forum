package main

import (
	"fmt"
	// db "reboot01.com/js/forum/database/sql-utils"
	"log"
	"net/http"

	s "reboot01.com/js/forum/src/server"
)

func init() {
	// db.DropDataBase()
	// db.DataBase()
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", s.ReverseMiddleware(s.LoginPage))
	mux.HandleFunc("POST /logout", s.AuthMiddleware(s.Logout))
	mux.HandleFunc("/signup", s.ReverseMiddleware(s.SignupPage))
	mux.HandleFunc("GET /home", s.AuthMiddleware(s.HomePage))
	mux.HandleFunc("/create", s.AuthMiddleware(s.CreatePage))
	mux.HandleFunc("GET /notifications", s.AuthMiddleware(s.NotificationsPage))
	mux.HandleFunc("GET /post/{postID}", s.AuthMiddleware(s.PostPage))
	mux.HandleFunc("POST /deletepost", s.AuthMiddleware(s.DeletePost))
	mux.HandleFunc("POST /addcomment", s.AuthMiddleware(s.AddComment))

	fmt.Println("Server running on http://localhost:8080\nTo stop the server press Ctrl+C")

	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
