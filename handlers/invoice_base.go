package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceHandler struct {
	DB *mongo.Database
}

func NewInvoiceHandler(db *mongo.Database) *InvoiceHandler {
	return &InvoiceHandler{
		DB: db,
	}
}
