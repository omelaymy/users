package api

import "github.com/google/uuid"

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Admin    bool   `json:"admin"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Admin    bool      `json:"admin"`
}

type UserIdResponse struct {
	Id uuid.UUID `json:"id"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
