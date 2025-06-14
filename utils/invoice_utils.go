package utils

import (
	"InvoicelyX/models"
)

// CalculateInvoiceTotal calculates the total amount
func CalculateInvoiceTotal(items []models.Item) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// ValidateInvoiceItems validates that all items have valid data
func ValidateInvoiceItems(items []models.Item) bool {
	if len(items) == 0 {
		return false
	}

	for _, item := range items {
		if item.Name == "" || item.Price <= 0 || item.Quantity <= 0 {
			return false
		}
	}
	return true
}
