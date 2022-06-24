package models

import "time"

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt time.Time o nulo
	UpdatedAt time.Time `json:"updatedAt"`
}

type ParamsQuery struct {
	Limit  int
	Offset int
}
