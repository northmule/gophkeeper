package models

import "time"

// User пользователь
type User struct {
	Common
	Login            string    `json:"login"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	CreatedAt        time.Time `json:"created_at"`
	PublicKey        string    `json:"public_key"`
	PrivateClientKey string    `json:"private_client_key"`
}
