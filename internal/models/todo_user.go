package models

import "time"

type User struct {
	ID         string    `json:"id" db:"id"`
	Mail       string    `json:"mail" db:"email"`
	Password   string    `json:"-" db:"password"`
	Created_at time.Time `json:"created_at" db:"created_at"`
	Updated_at time.Time `json:"updated_at" db:"updated_at"`
}
