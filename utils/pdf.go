package utils

import (
	"InvoicelyX/models"
	"bytes"
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

// GenerateInvoicePDF generates a PDF invoice from the given invoice data
func GenerateInvoicePDF(invoice *models.Invoice) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font for title
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 0, 0)

	// Title
	pdf.CellFormat(190, 15, "INVOICE", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Company/Invoice Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Invoice ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(80, 8, invoice.InvoiceID)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(80, 8, invoice.CreatedAt.Format("2006-01-02"))
	pdf.Ln(15)

	// Customer Information
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 8, "Bill To:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, "Customer:")
	pdf.Cell(80, 8, invoice.Customer)
	pdf.Ln(8)

	pdf.Cell(40, 8, "Email:")
	pdf.Cell(80, 8, invoice.Email)
	pdf.Ln(15)

	// Items Table Header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(60, 10, "Item Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Unit Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 1, "C", true, 0, "")

	// Items Table Content
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)

	for _, item := range invoice.Items {
		itemTotal := float64(item.Quantity) * item.Price

		pdf.CellFormat(60, 8, item.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, strconv.Itoa(item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("$%.2f", item.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("$%.2f", itemTotal), "1", 1, "R", false, 0, "")
	}

	// Total
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(130, 10, "TOTAL:", "0", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("$%.2f", invoice.Total), "1", 1, "R", true, 0, "")

	// Footer
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 8, "Thank you for your business!")
	// Get PDF as byte array
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
