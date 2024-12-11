package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UUID      string    `json:"uuid"`
}
