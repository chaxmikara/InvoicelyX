package utils_test

import (
	"InvoicelyX/models"
	"InvoicelyX/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCalculateInvoiceTotal tests the invoice total calculation function
func TestCalculateInvoiceTotal(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected float64
	}{
		{
			name: "Single item",
			items: []models.Item{
				{Name: "Item 1", Price: 10.50, Quantity: 2},
			},
			expected: 21.0,
		},
		{
			name: "Multiple items",
			items: []models.Item{
				{Name: "Item 1", Price: 10.50, Quantity: 2},
				{Name: "Item 2", Price: 5.25, Quantity: 4},
				{Name: "Item 3", Price: 100.0, Quantity: 1},
			},
			expected: 142.0, // (10.50 * 2) + (5.25 * 4) + (100.0 * 1) = 21 + 21 + 100 = 142
		},
		{
			name:     "Empty items",
			items:    []models.Item{},
			expected: 0.0,
		},
		{
			name: "Items with decimal prices",
			items: []models.Item{
				{Name: "Item 1", Price: 9.99, Quantity: 3},
				{Name: "Item 2", Price: 0.01, Quantity: 100},
			},
			expected: 30.97, // (9.99 * 3) + (0.01 * 100) = 29.97 + 1.00 = 30.97
		},
		{
			name: "Large quantities",
			items: []models.Item{
				{Name: "Item 1", Price: 1.0, Quantity: 1000},
			},
			expected: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CalculateInvoiceTotal(tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidateInvoiceItems tests the invoice items validation function
func TestValidateInvoiceItems(t *testing.T) {
	tests := []struct {
		name     string
		items    []models.Item
		expected bool
	}{
		{
			name: "Valid items",
			items: []models.Item{
				{Name: "Valid Item", Price: 10.0, Quantity: 1},
			},
			expected: true,
		},
		{
			name: "Multiple valid items",
			items: []models.Item{
				{Name: "Item 1", Price: 10.0, Quantity: 1},
				{Name: "Item 2", Price: 20.0, Quantity: 2},
			},
			expected: true,
		},
		{
			name:     "Empty items array",
			items:    []models.Item{},
			expected: false,
		},
		{
			name: "Item with empty name",
			items: []models.Item{
				{Name: "", Price: 10.0, Quantity: 1},
			},
			expected: false,
		},
		{
			name: "Item with zero price",
			items: []models.Item{
				{Name: "Invalid Item", Price: 0.0, Quantity: 1},
			},
			expected: false,
		},
		{
			name: "Item with negative price",
			items: []models.Item{
				{Name: "Invalid Item", Price: -5.0, Quantity: 1},
			},
			expected: false,
		},
		{
			name: "Item with zero quantity",
			items: []models.Item{
				{Name: "Invalid Item", Price: 10.0, Quantity: 0},
			},
			expected: false,
		},
		{
			name: "Item with negative quantity",
			items: []models.Item{
				{Name: "Invalid Item", Price: 10.0, Quantity: -1},
			},
			expected: false,
		},
		{
			name: "Mixed valid and invalid items",
			items: []models.Item{
				{Name: "Valid Item", Price: 10.0, Quantity: 1},
				{Name: "", Price: 20.0, Quantity: 1}, // Invalid name
			},
			expected: false,
		},
		{
			name: "Item with very small price",
			items: []models.Item{
				{Name: "Cheap Item", Price: 0.01, Quantity: 1},
			},
			expected: true,
		},
		{
			name: "Item with large price and quantity",
			items: []models.Item{
				{Name: "Expensive Item", Price: 999999.99, Quantity: 1000},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ValidateInvoiceItems(tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests for performance
func BenchmarkCalculateInvoiceTotal(b *testing.B) {
	items := []models.Item{
		{Name: "Item 1", Price: 10.50, Quantity: 2},
		{Name: "Item 2", Price: 5.25, Quantity: 4},
		{Name: "Item 3", Price: 100.0, Quantity: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.CalculateInvoiceTotal(items)
	}
}

func BenchmarkValidateInvoiceItems(b *testing.B) {
	items := []models.Item{
		{Name: "Item 1", Price: 10.50, Quantity: 2},
		{Name: "Item 2", Price: 5.25, Quantity: 4},
		{Name: "Item 3", Price: 100.0, Quantity: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.ValidateInvoiceItems(items)
	}
}
