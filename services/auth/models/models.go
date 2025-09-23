package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=4,max=20"`
	Email           string `json:"email" validate:"required,email,max=30"`
	Password        string `json:"password" validate:"required,min=4,max=10"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=30"`
	Password string `json:"password" validate:"required,min=4,max=10"`
}

type LoginHistory struct {
	ID        int       `json:"id"`
	UserEmail string    `json:"user_email"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
