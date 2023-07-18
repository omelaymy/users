package users

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Admin    bool      `json:"admin"`
	Password string    `json:"password,omitempty"`
}
