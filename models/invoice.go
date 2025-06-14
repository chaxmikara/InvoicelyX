package models

import (
	"time"
)

type Item struct {
	Name     string  `bson:"name" json:"name" validate:"required"`
	Price    float64 `bson:"price" json:"price" validate:"required,gt=0"`
	Quantity int     `bson:"quantity" json:"quantity" validate:"required,gt=0"`
}

// Invoice represents the complete invoice structure
type Invoice struct {
	InvoiceID string    `bson:"invoice_id" json:"invoice_id"`
	UserID    string    `bson:"user_id" json:"user_id"` // Associate invoice with user
	Customer  string    `bson:"customer" json:"customer" validate:"required,min=2,max=100"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Items     []Item    `bson:"items" json:"items" validate:"required,dive"`
	Total     float64   `bson:"total" json:"total"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
