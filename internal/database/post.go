package database

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Post struct {
	UserID     int            `db:"user_id"`
	Username   string         `db:"username"`
	FirstName  string         `db:"first_name"`
	LastName   string         `db:"last_name"`
	Likes      int            `db:"likes"`
	Dislikes   int            `db:"dislikes"`
	Comments   int            `db:"comments"`
	Categories []Category     `db:"categories"`
}

func (db *DB) GetAllCategories() ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var categories []Category

	query := `SELECT id, name FROM category`

	err := db.GetContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (db *DB) GetCategoriesForPost(postID int) ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var categories []Category

	query := `SELECT c.category_id, c.name 
              FROM category c
              JOIN post_category pc ON c.category_id = pc.category_id
              WHERE pc.post_id = $1`

	err := db.GetContext(ctx, &categories, query, postID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (db *DB) GetPostByID(postID int) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var post Post

	query := `SELECT p.user_id, p.username, p.first_name, p.last_name, p.likes, p.dislikes, p.comments
              FROM post p
              WHERE p.id = $1`

	err := db.GetContext(ctx, &post, query, postID)
	if err != nil {
		return nil, err
	}

	post.Categories, err = db.GetCategoriesForPost(postID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (db *DB) GetAllPosts() ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var posts []Post

	query := `SELECT p.user_id, p.username, p.first_name, p.last_name, p.avatar, p.likes, p.dislikes, p.comments
              FROM post p`

	err := db.GetContext(ctx, &posts, query)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		posts[i].Categories, err = db.GetCategoriesForPost(posts[i].UserID)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (db *DB) InsertPost(content string, image sql.NullString, userID string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()

    query := "INSERT INTO post (image, content, post_at, user_id) VALUES (?, ?, ?, ?)"
    stmt, err := db.PrepareContext(ctx, query)
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

    res, err := stmt.ExecContext(ctx, image, content, time.Now(), userID)
    if err != nil {
        return 0, err
    }

    lastID, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(lastID), nil
}

func (db *DB) InsertPostCategory(postID int, categoryID int) error {
    ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
    defer cancel()

    query := "INSERT INTO post_has_category (post_id, category_id) VALUES (?, ?)"
    stmt, err := db.PrepareContext(ctx, query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.ExecContext(ctx, postID, categoryID)
    if err != nil {
        return err
    }

    return nil
}