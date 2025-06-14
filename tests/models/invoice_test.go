package models_test

import (
	"InvoicelyX/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItemModel(t *testing.T) {
	t.Run("Valid Item Creation", func(t *testing.T) {
		item := models.Item{
			Name:     "Test Product",
			Price:    19.99,
			Quantity: 5,
		}

		assert.Equal(t, "Test Product", item.Name)
		assert.Equal(t, 19.99, item.Price)
		assert.Equal(t, 5, item.Quantity)
	})

	t.Run("Item with Zero Values", func(t *testing.T) {
		item := models.Item{}

		assert.Empty(t, item.Name)
		assert.Equal(t, 0.0, item.Price)
		assert.Equal(t, 0, item.Quantity)
	})

	t.Run("Item with High Values", func(t *testing.T) {
		item := models.Item{
			Name:     "Expensive Product",
			Price:    999999.99,
			Quantity: 1000,
		}

		assert.Equal(t, "Expensive Product", item.Name)
		assert.Equal(t, 999999.99, item.Price)
		assert.Equal(t, 1000, item.Quantity)
	})

	t.Run("Item with Decimal Quantities", func(t *testing.T) {
		// Note: Quantity is int, so this tests behavior with decimal prices
		item := models.Item{
			Name:     "Service Item",
			Price:    0.01, // Very small price
			Quantity: 1,
		}

		assert.Equal(t, "Service Item", item.Name)
		assert.Equal(t, 0.01, item.Price)
		assert.Equal(t, 1, item.Quantity)
	})

	t.Run("Item with Special Characters", func(t *testing.T) {
		item := models.Item{
			Name:     "Special Item - 50% off (was $100.00)",
			Price:    50.00,
			Quantity: 2,
		}

		assert.Equal(t, "Special Item - 50% off (was $100.00)", item.Name)
		assert.Equal(t, 50.00, item.Price)
		assert.Equal(t, 2, item.Quantity)
	})
}

func TestInvoiceModel(t *testing.T) {
	t.Run("Valid Invoice Creation", func(t *testing.T) {
		now := time.Now()
		items := []models.Item{
			{Name: "Product A", Price: 10.50, Quantity: 2},
			{Name: "Product B", Price: 25.00, Quantity: 1},
		}

		invoice := models.Invoice{
			InvoiceID: "inv-123",
			UserID:    "user-456",
			Customer:  "ABC Company",
			Email:     "billing@abc.com",
			Items:     items,
			Total:     46.00, // (10.50 * 2) + (25.00 * 1)
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, "inv-123", invoice.InvoiceID)
		assert.Equal(t, "user-456", invoice.UserID)
		assert.Equal(t, "ABC Company", invoice.Customer)
		assert.Equal(t, "billing@abc.com", invoice.Email)
		assert.Len(t, invoice.Items, 2)
		assert.Equal(t, 46.00, invoice.Total)
		assert.False(t, invoice.CreatedAt.IsZero())
		assert.False(t, invoice.UpdatedAt.IsZero())
	})

	t.Run("Invoice with Empty Items", func(t *testing.T) {
		invoice := models.Invoice{
			InvoiceID: "inv-empty",
			UserID:    "user-123",
			Customer:  "Test Customer",
			Email:     "test@example.com",
			Items:     []models.Item{},
			Total:     0.0,
		}

		assert.Equal(t, "inv-empty", invoice.InvoiceID)
		assert.Empty(t, invoice.Items)
		assert.Equal(t, 0.0, invoice.Total)
	})

	t.Run("Invoice with Single Item", func(t *testing.T) {
		item := models.Item{Name: "Single Product", Price: 99.99, Quantity: 1}
		invoice := models.Invoice{
			InvoiceID: "inv-single",
			UserID:    "user-789",
			Customer:  "Single Customer",
			Email:     "single@example.com",
			Items:     []models.Item{item},
			Total:     99.99,
		}

		assert.Equal(t, "inv-single", invoice.InvoiceID)
		assert.Len(t, invoice.Items, 1)
		assert.Equal(t, "Single Product", invoice.Items[0].Name)
		assert.Equal(t, 99.99, invoice.Total)
	})

	t.Run("Invoice with Multiple Items", func(t *testing.T) {
		items := []models.Item{
			{Name: "Item 1", Price: 10.00, Quantity: 3},
			{Name: "Item 2", Price: 15.50, Quantity: 2},
			{Name: "Item 3", Price: 5.25, Quantity: 10},
		}

		invoice := models.Invoice{
			InvoiceID: "inv-multi",
			UserID:    "user-multi",
			Customer:  "Multi Item Customer",
			Email:     "multi@example.com",
			Items:     items,
			Total:     113.50, // (10*3) + (15.50*2) + (5.25*10) = 30 + 31 + 52.50
		}

		assert.Equal(t, "inv-multi", invoice.InvoiceID)
		assert.Len(t, invoice.Items, 3)
		assert.Equal(t, 113.50, invoice.Total)

		// Verify each item
		assert.Equal(t, "Item 1", invoice.Items[0].Name)
		assert.Equal(t, 10.00, invoice.Items[0].Price)
		assert.Equal(t, 3, invoice.Items[0].Quantity)

		assert.Equal(t, "Item 2", invoice.Items[1].Name)
		assert.Equal(t, 15.50, invoice.Items[1].Price)
		assert.Equal(t, 2, invoice.Items[1].Quantity)

		assert.Equal(t, "Item 3", invoice.Items[2].Name)
		assert.Equal(t, 5.25, invoice.Items[2].Price)
		assert.Equal(t, 10, invoice.Items[2].Quantity)
	})

	t.Run("Invoice Default Values", func(t *testing.T) {
		invoice := models.Invoice{}

		assert.Empty(t, invoice.InvoiceID)
		assert.Empty(t, invoice.UserID)
		assert.Empty(t, invoice.Customer)
		assert.Empty(t, invoice.Email)
		assert.Nil(t, invoice.Items)
		assert.Equal(t, 0.0, invoice.Total)
		assert.True(t, invoice.CreatedAt.IsZero())
		assert.True(t, invoice.UpdatedAt.IsZero())
	})

	t.Run("Invoice with Long Customer Name", func(t *testing.T) {
		longCustomerName := "Very Long Customer Name That Exceeds Normal Limits And Contains Many Words To Test Edge Cases"

		invoice := models.Invoice{
			InvoiceID: "inv-long",
			UserID:    "user-long",
			Customer:  longCustomerName,
			Email:     "long@example.com",
			Items:     []models.Item{{Name: "Test", Price: 1.0, Quantity: 1}},
			Total:     1.0,
		}

		assert.Equal(t, longCustomerName, invoice.Customer)
		assert.True(t, len(invoice.Customer) > 50) // Testing edge case
	})

	t.Run("Invoice with Special Email Formats", func(t *testing.T) {
		specialEmails := []string{
			"test+tag@example.com",
			"user.name@sub.domain.com",
			"test_user@example-company.co.uk",
			"123numbers@example.org",
		}

		for i, email := range specialEmails {
			invoice := models.Invoice{
				InvoiceID: "inv-special-" + string(rune(i)),
				UserID:    "user-special",
				Customer:  "Special Customer",
				Email:     email,
				Items:     []models.Item{{Name: "Test", Price: 1.0, Quantity: 1}},
				Total:     1.0,
			}

			assert.Equal(t, email, invoice.Email)
		}
	})
}

func TestInvoiceItemsValidation(t *testing.T) {
	t.Run("Valid Items Slice", func(t *testing.T) {
		items := []models.Item{
			{Name: "Valid Item 1", Price: 10.0, Quantity: 1},
			{Name: "Valid Item 2", Price: 20.0, Quantity: 2},
		}

		// Test that items can be created and accessed
		assert.Len(t, items, 2)
		assert.Equal(t, "Valid Item 1", items[0].Name)
		assert.Equal(t, "Valid Item 2", items[1].Name)
	})

	t.Run("Items with Edge Case Values", func(t *testing.T) {
		items := []models.Item{
			{Name: "Minimum Price", Price: 0.01, Quantity: 1},
			{Name: "Large Quantity", Price: 1.0, Quantity: 10000},
			{Name: "High Price", Price: 999999.99, Quantity: 1},
		}

		assert.Len(t, items, 3)
		assert.Equal(t, 0.01, items[0].Price)
		assert.Equal(t, 10000, items[1].Quantity)
		assert.Equal(t, 999999.99, items[2].Price)
	})
}

// Benchmark tests
func BenchmarkItemCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := models.Item{
			Name:     "Benchmark Product",
			Price:    19.99,
			Quantity: 5,
		}
		_ = item
	}
}

func BenchmarkInvoiceCreation(b *testing.B) {
	items := []models.Item{
		{Name: "Product A", Price: 10.50, Quantity: 2},
		{Name: "Product B", Price: 25.00, Quantity: 1},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		invoice := models.Invoice{
			InvoiceID: "inv-123",
			UserID:    "user-456",
			Customer:  "ABC Company",
			Email:     "billing@abc.com",
			Items:     items,
			Total:     46.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_ = invoice
	}
}

func BenchmarkItemsSliceCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items := []models.Item{
			{Name: "Item 1", Price: 10.00, Quantity: 1},
			{Name: "Item 2", Price: 20.00, Quantity: 2},
			{Name: "Item 3", Price: 30.00, Quantity: 3},
		}
		_ = items
	}
}
