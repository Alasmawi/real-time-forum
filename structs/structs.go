package structs

import (
	"database/sql"
	"log"
	"html/template"
	"path/filepath"
	"time"
)

var Db *sql.DB
var Templates *template.Template

func init() {
	Templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))

	var err error
	Db, err = sql.Open("sqlite3", "./database/main.db")
	if err != nil {
		log.Println("Database connection failed")
		log.Fatal(err)
	}
	// var err error
	// Db, err = sql.Open("sqlite3", "./database/main.db")
	// if err != nil {
	// 	log.Println("Database connection failed")
	// 	log.Fatal(err)
	// }
	// return
}

type ErrorPageData struct {
	Code     int
	ErrorMsg string
}

type PageData struct {
	HasSession      bool
	UserID          int
	UserName        string
	Avatar          string
	RoleName        string
	TotalLikes      int
	TotalPosts      int
	Categories      []Category
	Users           []User
	Posts           []Post
	TotalUsers      int
	TotalCategories int
	Notifications   []Notification
	RoleID          int
	Post            Post
	Comments        []Comment
}

type User struct {
	ID               int            `json:"id"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Username         string         `json:"username"`
	Email            string         `json:"email"`
	Password         string         `json:"password"`
	SessionSessionID int            `json:"current_session"`
	RoleID           int            `json:"role_id"`
	Avatar           sql.NullString `json:"avatar"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	FirstName string
	LastName  string
	Username  string
	Content   string
	CreatedAt time.Time
	Avatar    sql.NullString // Add this field
	Likes     int
	Dislikes  int
}

type Post struct {
	PostID     int
	Image      sql.NullString
	Content    string
	PostAt     time.Time
	UserUserID int
	Username   string
	FirstName  string
	LastName   string
	Avatar     sql.NullString
	Likes      int
	Dislikes   int
	Comments   int
	Categories []Category
}

type Notification struct {
	ID        int
	UserID    int
	PostID    int
	Message   string
	CreatedAt time.Time
	UserImage string
	UserName  string
}
