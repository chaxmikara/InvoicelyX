package utils

import (
	"InvoicelyX/models"
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

// generates a CSV file from the given invoices data
func GenerateInvoicesCSV(invoices []models.Invoice) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV headers
	headers := []string{
		"Invoice ID",
		"Customer",
		"Email",
		"Items",
		"Total",
		"Created At",
	}

	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("error writing CSV headers: %v", err)
	}

	// Write invoice data
	for _, invoice := range invoices {
		// Format items as a string
		var itemStrings []string
		for _, item := range invoice.Items {
			itemStr := fmt.Sprintf("%s (Qty: %d, Price: $%.2f)",
				item.Name, item.Quantity, item.Price)
			itemStrings = append(itemStrings, itemStr)
		}
		itemsStr := strings.Join(itemStrings, "; ")

		record := []string{
			invoice.Customer,
			invoice.Email,
			itemsStr,
			fmt.Sprintf("$%.2f", invoice.Total),
			invoice.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("error writing CSV record: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return buf.Bytes(), nil
}
