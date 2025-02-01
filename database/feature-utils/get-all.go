package database

import (
	"database/sql"
	"log"

	strct "reboot01.com/js/forum/structs"
)

func GetAllCategories() ([]strct.Category, error) {
	rows, err := strct.Db.Query("SELECT category_id, name, description FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []strct.Category
	for rows.Next() {
		var category strct.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetAllUsers() ([]strct.User, error) {
	rows, err := strct.Db.Query("SELECT user_id, f_name, l_name, username, email, avatar FROM user")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var users []strct.User
	for rows.Next() {
		var user strct.User
		var avatar sql.NullString
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &avatar); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		user.Avatar = avatar
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return users, nil
}
