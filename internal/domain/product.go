package domain

import (
	"time"
)

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

func (p *Product) Validate() error {
	if p.Name == "" || p.Price <= 0 || p.Stock < 0 || p.CategoryID == "" {
		return ErrInvalidProduct
	}
	return nil
}
