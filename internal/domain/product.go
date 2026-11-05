package domain

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  string
	Specs       map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
