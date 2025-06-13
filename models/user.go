package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string    `bson:"user_id" json:"user_id,omitempty"`
	FirstName string    `bson:"firstName" json:"firstName" validate:"required,min=2,max=50"`
	LastName  string    `bson:"lastName" json:"lastName" validate:"required,min=2,max=50"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Password  string    `bson:"password" json:"password,omitempty" validate:"required,min=8,max=100"`
	Role      string    `bson:"role" json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	IsActive  bool      `bson:"isActive" json:"isActive,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LogoutRequest struct {
	SessionToken string `json:"sessionToken" validate:"required"`
}
