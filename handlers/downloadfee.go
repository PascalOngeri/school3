package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	
)

// GenerateFeeStatement generates and downloads the fee statement as a PDF
func GenerateFeeStatement(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	adm := r.FormValue("adm")
	if adm == "" {
		http.Error(w, "Admission number is required", http.StatusBadRequest)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)

	// Add logo
	logoPath := filepath.Join("static", "logo.png")
	pdf.Image(logoPath, 80, 10, 50, 0, false, "", 0, "")

	// Title
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Fee Statement", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Admission Number: %s", adm), "", 1, "C", false, 0, "")

	// Table Headers
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Receipt No.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Date", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Balance", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	// Fetch data from DB
	rows, err := db.Query("SELECT id, date, amount, bal FROM payment WHERE adm = ? ORDER BY id ASC", adm)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pdf.SetFont("Arial", "", 10)
	for rows.Next() {
		var id, date, amount, balance string
		if err := rows.Scan(&id, &date, &amount, &balance); err != nil {
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		pdf.CellFormat(40, 10, id, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, date, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, amount, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, balance, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	// Output the PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=fee_statement.pdf")
	_ = pdf.Output(w)
}

// GenerateFeeStructure generates and downloads the fee structure as a PDF
func GenerateFeeStructure(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	class := r.FormValue("genclass")
	if class == "" {
		http.Error(w, "Class is required", http.StatusBadRequest)
		return
	}

	// PDF Setup
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)

	// Add logo
	logoPath := filepath.Join("static", "logo.png")
	pdf.Image(logoPath, 80, 10, 50, 0, false, "", 0, "")

	// Title
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Fee Structure", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Class: %s", class), "", 1, "C", false, 0, "")

	// Table Headers
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"S.No", "Payment Name", "Term 1", "Term 2", "Term 3", "Total"}
	for _, header := range headers {
		pdf.CellFormat(38, 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Fetch data from DB
	rows, err := db.Query("SELECT id, paymentname, term1, term2, term3, amount FROM feepay WHERE form = ?", class)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	pdf.SetFont("Arial", "", 10)
	for rows.Next() {
		var id, paymentName, term1, term2, term3, total string
		if err := rows.Scan(&id, &paymentName, &term1, &term2, &term3, &total); err != nil {
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		pdf.CellFormat(38, 10, id, "1", 0, "", false, 0, "")
		pdf.CellFormat(38, 10, paymentName, "1", 0, "", false, 0, "")
		pdf.CellFormat(38, 10, term1, "1", 0, "", false, 0, "")
		pdf.CellFormat(38, 10, term2, "1", 0, "", false, 0, "")
		pdf.CellFormat(38, 10, term3, "1", 0, "", false, 0, "")
		pdf.CellFormat(38, 10, total, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	// Output the PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=fee_structure.pdf")
	_ = pdf.Output(w)
}
