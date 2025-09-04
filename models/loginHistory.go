package models

import "time"

type LoginHistory struct {
	ID        int       `json:"id"`
	UserEmail string    `json:"user_email"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
