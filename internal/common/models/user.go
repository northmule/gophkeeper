package models

import "time"

type User struct {
	Common
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
