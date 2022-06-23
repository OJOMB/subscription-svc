package domain

import "time"

type User struct {
	ID          string     `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	DateOfBirth *time.Time `json:"dob"`
}
