package main

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Post struct {
	ID     int    `db:"id" json:"id"`
	UserID int    `db:"user_id" json:"user_id"`
	Text   string `db:"text" json:"text"`
}

func InitDB() *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", "app.db")

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		text TEXT
	);
	`

	db.MustExec((schema))
	return db
}

func GetUserHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []User
		err := db.Select(&users, "SELECT * FROM users")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

func CreateUserHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		res, err := db.Exec("INSERT INTO users (name) VALUES (?)", user.Name)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		id, _ := res.LastInsertId()
		user.ID = int(id)

		json.NewEncoder(w).Encode(user)
	}
}

func GetPostHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var posts []Post
		err := db.Select(&posts, "SELECT * FROM posts")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(posts)
	}
}

func CreatePostHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		res, err := db.Exec("INSERT INTO posts (user_id,text) VALUES (?,?)", post.UserID, post.Text)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		id, _ := res.LastInsertId()
		post.ID = int(id)

		json.NewEncoder(w).Encode(post)
	}
}
