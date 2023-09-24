package models

import (
	"database/sql"
	"time"

    "github.com/go-chi/jwtauth"
)

var DB *sql.DB
var TokenAuth *jwtauth.JWTAuth

type User struct {
	FirstName string `json:"firstName"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt time.Time
    EditedAt time.Time 
    Organisation []string
}


