package main

import (
	"fmt"
	db "forum/database/sql-utils"
	"forum/src/server"
	"log"
	"net/http"
)

func init() {
	db.DataBase()
	//db.DropDataBase()
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// http.HandleFunc("/", server.MainPage)
	http.HandleFunc("/", server.ReverseMiddleware(server.LoginPage))
	http.HandleFunc("/logout", server.AuthMiddleware(server.Logout))
	http.HandleFunc("/signup", server.SignupPage)
	http.HandleFunc("/home", server.HomePage)
	http.HandleFunc("/newpost", server.AuthMiddleware(server.NewPostPage))
	http.HandleFunc("/notifications", server.AuthMiddleware(server.NotificationsPage))
	http.HandleFunc("/post", server.AuthMiddleware(server.PostPage))
	http.HandleFunc("/deletepost", server.AuthMiddleware(server.DeletePost))
	http.HandleFunc("/addcomment", server.AuthMiddleware(server.AddComment))

	fmt.Println("Server running on http://localhost:8080\nTo stop the server press Ctrl+C")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
