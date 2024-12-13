package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"bytes"


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
	Type string
	Term string

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
	for rows.Next() {
		var payment PaymentRecord
		if err := rows.Scan(&payment.ID, &payment.Date, &payment.Amount, &payment.Balance); err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			continue
		}
		payments = append(payments, payment)
	}
	if rows.Err() != nil {
		log.Printf("Rows iteration error: %v\n", rows.Err())
		http.Error(w, "Error processing payment records", http.StatusInternalServerError)
		return
	}

	// Fetch school name and logo
	var schoolName, icon string
	err = db.QueryRow("SELECT name, icon FROM api").Scan(&schoolName, &icon)
	if err != nil {
		log.Printf("Error fetching school name: %v\n", err)
		http.Error(w, "Error fetching school name", http.StatusInternalServerError)
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add logo and school name
	pdf.ImageOptions(icon, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.Ln(50)
	pdf.CellFormat(0, 10, schoolName, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, "Admission Number "+adm+" Fee Statement", "", 1, "C", false, 0, "")

	// Add table headers
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Payment No.", "Date", "Amount", "Balance", "Status"}
	widths := []float64{40, 38, 38, 38, 38}
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Add table rows
	pdf.SetFont("Arial", "", 10)
	for _, payment := range payments {
		status := "Received"
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", payment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, payment.Date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", payment.Amount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", payment.Balance), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, status, "1", 1, "C", false, 0, "")
	}

	// Generate PDF to memory
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating PDF: %v\n", err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Set headers and write PDF content to response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=Statement.pdf")
	w.Write(buf.Bytes())
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

	// Query the database for optional payments from otherpay table
	optionalRows, err := db.Query("SELECT id, t1, t2, t3, amount,type FROM other ORDER BY id ASC")
	if err != nil {
		log.Printf("Database query error for optional payments: %v\n", err)
		http.Error(w, "Error fetching optional payments", http.StatusInternalServerError)
		return
	}
	defer optionalRows.Close()

	// Create a slice for optional payments
	var optionalPayments []PaymentRecord
	for optionalRows.Next() {
		var optionalPayment PaymentRecord
		if err := optionalRows.Scan(&optionalPayment.ID, &optionalPayment.Term1, &optionalPayment.Term2, &optionalPayment.Term3, &optionalPayment.Amount,&optionalPayment.Type); err != nil {
			log.Printf("Failed to scan optional payment row: %v\n", err)
			continue
		}
		optionalPayments = append(optionalPayments, optionalPayment)
	}

	// Check for any row iteration errors for optional payments
	if optionalRows.Err() != nil {
		log.Printf("Optional payments iteration error: %v\n", optionalRows.Err())
		http.Error(w, "Error processing optional payments", http.StatusInternalServerError)
		return
	}

	// Query the database for transportation payments from the bus table
	transportationRows, err := db.Query("SELECT area, t1, t2, t3, amount, id FROM bus ORDER BY id ASC")
	if err != nil {
		log.Printf("Database query error for transportation payments: %v\n", err)
		http.Error(w, "Error fetching transportation payments", http.StatusInternalServerError)
		return
	}
	defer transportationRows.Close()

	// Create a slice for transportation payments
	var transportationPayments []PaymentRecord
	for transportationRows.Next() {
		var transportPayment PaymentRecord
		if err := transportationRows.Scan(&transportPayment.PaymentName, &transportPayment.Term1, &transportPayment.Term2, &transportPayment.Term3, &transportPayment.Amount, &transportPayment.ID); err != nil {
			log.Printf("Failed to scan transportation payment row: %v\n", err)
			continue
		}
		transportationPayments = append(transportationPayments, transportPayment)
	}

	// Check for any row iteration errors for transportation payments
	if transportationRows.Err() != nil {
		log.Printf("Transportation payments iteration error: %v\n", transportationRows.Err())
		http.Error(w, "Error processing transportation payments", http.StatusInternalServerError)
		return
	}

// Remove the second declaration of 'err' by using assignment instead of short declaration

var schoolName, icon string
err = db.QueryRow("SELECT name, icon FROM api").Scan(&schoolName, &icon)
if err != nil {
    log.Printf("Error fetching school name: %v\n", err)
    http.Error(w, "Error fetching school name", http.StatusInternalServerError)
    return
}

// Create a new PDF document
pdf := gofpdf.New("P", "mm", "A4", "")
pdf.AddPage()

// Set logo path and add the logo to the document
pdf.ImageOptions(icon, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

// Add school name below the logo
pdf.SetFont("Arial", "B", 16)
pdf.Ln(50) // Adjust the vertical space as needed
pdf.CellFormat(0, 10, schoolName, "", 1, "C", false, 0, "") // Use schoolName from DB

	// Adjust the vertical space as needed
pdf.CellFormat(0, 10, adm + " Fee structure", "", 1, "C", false, 0, "")



	// Add Fee Payment Table
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Payment No.", "Payment Name", "Term 1", "Term 2", "Term 3", "Amount"}
	widths := []float64{30, 40, 30, 30, 30, 30} // Adjust the widths as needed
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Set font for table rows
	pdf.SetFont("Arial", "", 10)
	for _, payment := range payments {
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", payment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, payment.PaymentName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term2), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term3), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Amount), "1", 1, "C", false, 0, "")
	}

	// Optional Payment Table
	pdf.Ln(10) // Space before next table
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Optional Payments", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	optionalHeaders := []string{"Payment No.", "Payment Name", "Term 1", "Term2", "Term3", "Amount"}
	for i, header := range optionalHeaders {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, optionalPayment := range optionalPayments {
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", optionalPayment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, optionalPayment.Type, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", optionalPayment.Term1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", optionalPayment.Term2), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f",  optionalPayment.Term3), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", optionalPayment.Amount), "1", 1, "C", false, 0, "")
	}

	// Transportation Payment Table
	pdf.Ln(10) // Space before next table
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "Transportation Payments", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	transportationHeaders := []string{"Payment No.", "Payment Name", "Term1", "Term2", "Term3", "Amount"}
	for i, header := range transportationHeaders {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, transportPayment := range transportationPayments {
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", transportPayment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, transportPayment.PaymentName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", transportPayment.Term1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", transportPayment.Term2), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", transportPayment.Term3), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", transportPayment.Amount), "1", 1, "C", false, 0, "")
	}

	// Output the PDF to a file
	outputFileName := "fee_report.pdf"
	err = pdf.OutputFileAndClose(outputFileName)
	if err != nil {
		log.Printf("Error generating PDF: %v\n", err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Set the header to prompt a download of the generated PDF
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFileName)
	w.Header().Set("Content-Type", "application/pdf")

	// Serve the PDF file as response
	http.ServeFile(w, r, outputFileName)

	// Remove the file after serving
	os.Remove(outputFileName)
}

func Individualfee(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    // Retrieve cookies for admission number and class
    admCookie, err := r.Cookie("adm")
    if err != nil {
        log.Printf("Error getting adm cookie: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    cla, err := r.Cookie("form")
    if err != nil {
        log.Printf("Error getting form cookie: %v", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    adm := admCookie.Value
    class := cla.Value

    if adm == "" || class == "" {
        log.Println("Admission number and class are required")
        http.Error(w, "Admission number and class are required", http.StatusBadRequest)
        return
    }

    // Query the database for fee payment records
    rows, err := db.Query("SELECT id, paymentname, term1, term2, term3, amount FROM feepay WHERE form = ? ORDER BY id ASC", class)
    if err != nil {
        log.Printf("Database query error: %v\n", err)
        http.Error(w, "Error fetching payment records", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var payments []PaymentRecord
    var totalFee float64

    for rows.Next() {
        var payment PaymentRecord
        if err := rows.Scan(&payment.ID, &payment.PaymentName, &payment.Term1, &payment.Term2, &payment.Term3, &payment.Amount); err != nil {
            log.Printf("Failed to scan row: %v\n", err)
            continue
        }
        totalFee += payment.Amount
        payments = append(payments, payment)
    }

    if rows.Err() != nil {
        log.Printf("Rows iteration error: %v\n", rows.Err())
        http.Error(w, "Error processing payment records", http.StatusInternalServerError)
        return
    }

    // Query the database for optional payments
    optionalRows, err := db.Query("SELECT id, payname, term, amount, date FROM otherpayment WHERE adm=? ORDER BY id ASC", adm)
    if err != nil {
        log.Printf("Database query error for optional payments: %v\n", err)
        http.Error(w, "Error fetching optional payments", http.StatusInternalServerError)
        return
    }
    defer optionalRows.Close()

    var optionalPayments []PaymentRecord
    var totalOptionalFee float64

    for optionalRows.Next() {
        var optionalPayment PaymentRecord
        if err := optionalRows.Scan(&optionalPayment.ID, &optionalPayment.Type, &optionalPayment.Term, &optionalPayment.Amount, &optionalPayment.Date); err != nil {
            log.Printf("Failed to scan optional payment row: %v\n", err)
            continue
        }
        totalOptionalFee += optionalPayment.Amount
        optionalPayments = append(optionalPayments, optionalPayment)
    }

    if optionalRows.Err() != nil {
        log.Printf("Optional payments iteration error: %v\n", optionalRows.Err())
        http.Error(w, "Error processing optional payments", http.StatusInternalServerError)
        return
    }

    // Fetch school details
    var schoolName, icon string
    err = db.QueryRow("SELECT name, icon FROM api").Scan(&schoolName, &icon)
    if err != nil {
        log.Printf("Error fetching school name: %v\n", err)
        http.Error(w, "Error fetching school name", http.StatusInternalServerError)
        return
    }

    // Create the PDF document
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
   
    pdf.ImageOptions(icon, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

    pdf.SetFont("Arial", "B", 16)
    pdf.Ln(50)
    pdf.CellFormat(0, 10, schoolName, "", 1, "C", false, 0, "")
    pdf.CellFormat(0, 10, "Admission Number "+adm+" Fee Structure", "", 1, "C", false, 0, "")

    // Add Fee Payment Table
    pdf.SetFont("Arial", "B", 12)
    headers := []string{"Payment No.", "Payment Name", "Term 1", "Term 2", "Term 3", "Amount"}
    widths := []float64{30, 40, 30, 30, 30, 30}
    for i, header := range headers {
        pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)

    pdf.SetFont("Arial", "", 10)
    for _, payment := range payments {
        pdf.CellFormat(30, 10, fmt.Sprintf("%d", payment.ID), "1", 0, "C", false, 0, "")
        pdf.CellFormat(40, 10, payment.PaymentName, "1", 0, "C", false, 0, "")
        pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term1), "1", 0, "C", false, 0, "")
        pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term2), "1", 0, "C", false, 0, "")
        pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Term3), "1", 0, "C", false, 0, "")
        pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", payment.Amount), "1", 1, "C", false, 0, "")
    }
    pdf.SetFont("Arial", "B", 12)
    pdf.CellFormat(190, 10, fmt.Sprintf("Total Fee: %.2f", totalFee), "1", 1, "R", false, 0, "")

    // Add Optional Payments Table
    pdf.Ln(10)
    pdf.CellFormat(0, 10, "Optional Payments", "", 1, "L", false, 0, "")
    pdf.SetFont("Arial", "B", 12)
    optionalHeaders := []string{"Payment No.", "Payment Name", "Type", "Amount", "Date"}
    wi := []float64{30, 40, 40, 40, 40, }
    for i, header := range optionalHeaders {
        pdf.CellFormat(wi[i], 10, header, "1", 0, "C", false, 0, "")
    }
    pdf.Ln(-1)

    pdf.SetFont("Arial", "", 10)
    for _, optionalPayment := range optionalPayments {
        pdf.CellFormat(30, 10, fmt.Sprintf("%d", optionalPayment.ID), "1", 0, "C", false, 0, "")
        pdf.CellFormat(40, 10, optionalPayment.Type, "1", 0, "C", false, 0, "")
        pdf.CellFormat(40, 10, optionalPayment.Term, "1", 0, "C", false, 0, "")
        pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", optionalPayment.Amount), "1", 0, "C", false, 0, "")
        pdf.CellFormat(40, 10, optionalPayment.Date, "1", 1, "C", false, 0, "")
    }
    pdf.SetFont("Arial", "B", 12)
    pdf.CellFormat(190, 10, fmt.Sprintf("Total Optional Fee: %.2f", totalOptionalFee), "1", 1, "R", false, 0, "")

    // Output the PDF to a file
    outputFileName := "fee_report.pdf"
    err = pdf.OutputFileAndClose(outputFileName)
    if err != nil {
        log.Printf("Error generating PDF: %v\n", err)
        http.Error(w, "Error generating PDF", http.StatusInternalServerError)
        return
    }

    // Serve the PDF file as a response
    w.Header().Set("Content-Disposition", "attachment; filename="+outputFileName)
    w.Header().Set("Content-Type", "application/pdf")
    http.ServeFile(w, r, outputFileName)


	// Remove the file after serving
	os.Remove(outputFileName)
}


