package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PaymentRecord defines the structure for each payment record
type PaymentRecord struct {
	ID      int
	Date    string
	Amount  float64
	Balance float64
	PaymentName string
	Term1       float64
	Term2       float64
	Term3       float64

}

// GenerateFeeHandler generates the fee statement for a given admission number
func GenerateFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	adm := r.FormValue("adm")
	if adm == "" {
		log.Println("Admission number is required")
		http.Error(w, "Admission number is required", http.StatusBadRequest)
		return
	}
	
	// Query the database to get payment records
	rows, err := db.Query("SELECT id, date, amount, bal FROM payment WHERE adm = ? ORDER BY id ASC", adm)
	if err != nil {
		log.Printf("Database query error: %v\n", err)
		http.Error(w, "Error fetching payment records", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the payment records
	var payments []PaymentRecord

	// Populate the slice with the query results
	for rows.Next() {
		var payment PaymentRecord
		if err := rows.Scan(&payment.ID, &payment.Date, &payment.Amount, &payment.Balance); err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}
		payments = append(payments, payment)
	}

	// Check for any row iteration errors
	if rows.Err() != nil {
		log.Printf("Rows iteration error: %v\n", rows.Err())
		http.Error(w, "Error processing payment records", http.StatusInternalServerError)
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set logo path and add the logo to the document
	logoPath := filepath.Join("assets", "images", "logo.png")
	pdf.ImageOptions(logoPath, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	// Add school name below the logo
	pdf.SetFont("Arial", "B", 16)
	pdf.Ln(50) // Adjust the vertical space as needed

	pdf.CellFormat(0, 10, "INFINITY SCHOOLS", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add table headers
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Payment No.", "Date", "Amount", "Balance", "Status"}
	widths := []float64{40, 38, 38, 38, 38}
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Set font for table rows
	pdf.SetFont("Arial", "", 10)

	// Iterate over the payment records and add data to the PDF
	for _, payment := range payments {
		status := "Received"
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", payment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, payment.Date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", payment.Amount), "1", 0, "C", false, 0, "") // Display amount
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", payment.Balance), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, status, "1", 1, "C", false, 0, "")
	}

	// Check for any PDF generation errors
	if pdf.Err() {
		log.Println("Error generating PDF")
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Create output directory
	outputDir := "generated_pdfs"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory: %v", err)
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Generate unique file name
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("fee_statement_%s_%s.pdf", adm, timestamp)
	filePath := filepath.Join(outputDir, fileName)

	// Save the PDF to the file
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Printf("Error saving PDF: %v", err)
		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
		return
	}

	// Redirect to the parent page
	http.Redirect(w, r, "/parent", http.StatusSeeOther)
}
func GenerateFee(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the admission number from the form
	adm := r.FormValue("genclass")
	if adm == "" {
		log.Println("Admission number is required")
		http.Error(w, "Admission number is required", http.StatusBadRequest)
		return
	}

	// Query the database to get payment records for the given class
	rows, err := db.Query("SELECT id, paymentname, term1, term2, term3, amount FROM feepay WHERE form = ? ORDER BY id ASC", adm)
	if err != nil {
		log.Printf("Database query error: %v\n", err)
		http.Error(w, "Error fetching payment records", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the payment records
	var payments []PaymentRecord

	// Populate the slice with the query results
	for rows.Next() {
		var payment PaymentRecord
		if err := rows.Scan(&payment.ID, &payment.PaymentName, &payment.Term1, &payment.Term2, &payment.Term3, &payment.Amount); err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}
		payments = append(payments, payment)
	}

	// Check for any row iteration errors
	if rows.Err() != nil {
		log.Printf("Rows iteration error: %v\n", rows.Err())
		http.Error(w, "Error processing payment records", http.StatusInternalServerError)
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set logo path and add the logo to the document
	logoPath := filepath.Join("assets", "images", "logo.png")
	pdf.ImageOptions(logoPath, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	// Add school name below the logo
	pdf.SetFont("Arial", "B", 16)
	pdf.Ln(50) // Adjust the vertical space as needed
	pdf.CellFormat(0, 10, "INFINITY SCHOOLS", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add table headers
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Payment No.", "Payment Name", "Term 1", "Term 2", "Term 3", "Amount"}
	widths := []float64{30, 40, 30, 30, 30, 30} // Adjust the widths as needed
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Set font for table rows
	pdf.SetFont("Arial", "", 10)

	// Iterate over the payment records and add data to the PDF
	for _, payment := range payments {
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", payment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, payment.PaymentName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term2), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term3), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Amount), "1", 1, "C", false, 0, "")
	}

	// Check for any PDF generation errors
	if pdf.Err() {
		log.Println("Error generating PDF")
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Create output directory if it doesn't exist
	outputDir := "generated_pdfs"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory: %v", err)
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Generate unique file name based on admission number and timestamp
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("fee_statement_%s_%s.pdf", adm, timestamp)
	filePath := filepath.Join(outputDir, fileName)

	// Save the PDF to the file
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Printf("Error saving PDF: %v", err)
		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
		return
	}

	// Serve the generated PDF to the client for download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	http.ServeFile(w, r, filePath)
}
