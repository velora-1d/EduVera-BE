package pdf

import (
	"bytes"
	"fmt"
	"time"

	"prabogo/internal/model"

	"github.com/go-pdf/fpdf"
)

func GeneratePaySlipPDF(slip *model.PaySlip) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// -- Header --
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, slip.SchoolName, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	// Placeholder address if not available
	pdf.CellFormat(0, 5, "Jl. Pendidikan No. 1", "", 1, "C", false, 0, "")

	pdf.Ln(5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY()) // Horizontal line
	pdf.Ln(5)

	// -- Title --
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "SLIP GAJI KARYAWAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("Periode: %s", formatPeriod(slip.Payroll.Period)), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// -- Employee Info --
	leftX := 15.0
	pdf.SetX(leftX)
	pdf.SetFont("Arial", "", 11)

	// Helper for key-value row
	printRow := func(label, value string) {
		pdf.SetX(leftX)
		pdf.CellFormat(40, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 6, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6, value, "", 1, "L", false, 0, "")
	}

	printRow("Nama", slip.Employee.Name)
	printRow("NIP", slip.Employee.NIP) // TODO: Handle empty NIP safely if needed
	printRow("Jabatan", slip.Employee.Role)
	printRow("Status", slip.Employee.EmployeeType)

	pdf.Ln(10)

	// -- Details Table --
	// We'll split into Earnings (Penerimaan) and Deductions (Potongan)
	// But first, let's list all items

	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(100, 8, "Keterangan", "1", 0, "L", true, 0, "")
	pdf.CellFormat(90, 8, "Jumlah (Rp)", "1", 1, "R", true, 0, "")

	pdf.SetFont("Arial", "", 11)

	// Base Salary
	pdf.CellFormat(100, 7, "Gaji Pokok", "1", 0, "L", false, 0, "")
	pdf.CellFormat(90, 7, formatRupiah(slip.Payroll.BaseSalary), "1", 1, "R", false, 0, "")

	// Parse Details JSONB to display specific allowances/deductions
	// For now, assuming Details is available in the struct as []PayDetail or similar if mapped.
	// But mapped struct usually has it as JSONB (string/bytes) or interface{}.
	// If `slip.Payroll.Details` is raw JSON, we might rely on the summary Allowances/Deductions if detailed parsing is complex here without proper struct tag unmarshalling in model.
	// However, `slip.Payroll.Details` in model/sdm.go is `interface{}` or `datatypes.JSON`?
	// Checking `model/sdm.go` earlier: `Details JSONB \`json:"details\"\``.
	// In Go code it might be `pgtype.JSONB` or `json.RawMessage` or mapped struct depending on implementation.
	// Let's assume we map it to []model.PayDetail in Service before returning, OR we iterate `Allowances` and `Deductions` total if details are hard to parse here.

	// Since we defined `PayDetail` in service but didn't verify if it's in the `PaySlip` struct passed here.
	// The `PaySlip` struct in service.go was:
	// type PaySlip struct { Employee, Payroll, SchoolName }
	// And `Payroll` has `Details` which is mapped from DB.
	// Let's assume for MVP we show Total Allowances and Total Deductions if Details are not easily accessible as typed struct,
	// BUT better to check if we can cast it.
	// To be safe and since we don't want to import heavy JSON logic here if possible, let's use the totals for now,
	// unless we update the model to have TypedDetails.
	// Let's just list "Tunjangan Lain-lain" and "Potongan Lain-lain" using the totals.

	if slip.Payroll.Allowances > 0 {
		pdf.CellFormat(100, 7, "Total Tunjangan", "1", 0, "L", false, 0, "")
		pdf.CellFormat(90, 7, formatRupiah(slip.Payroll.Allowances), "1", 1, "R", false, 0, "")
	}

	if slip.Payroll.Deductions > 0 {
		pdf.CellFormat(100, 7, "Total Potongan", "1", 0, "L", false, 0, "")
		pdf.CellFormat(90, 7, fmt.Sprintf("- %s", formatRupiah(slip.Payroll.Deductions)), "1", 1, "R", false, 0, "")
	}

	// Net Salary
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(100, 10, "Total Diterima (Take Home Pay)", "1", 0, "L", false, 0, "")
	pdf.CellFormat(90, 10, formatRupiah(slip.Payroll.NetSalary), "1", 1, "R", false, 0, "")

	pdf.Ln(20)

	// -- Footer / Signature --
	// Left: Employee, Right: Treasurer/Principal
	currentY := pdf.GetY()

	// Employee
	pdf.SetX(15)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(60, 5, "Penerima,", "", 1, "C", false, 0, "")
	pdf.Ln(20)
	pdf.SetX(15)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(60, 5, slip.Employee.Name, "", 1, "C", false, 0, "")

	// Treasurer
	pdf.SetXY(130, currentY)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(60, 5, fmt.Sprintf("Karanganyar, %s", time.Now().Format("02 Jan 2006")), "", 1, "C", false, 0, "")
	pdf.SetXY(130, currentY+5)
	pdf.CellFormat(60, 5, "Bendahara Sekolah,", "", 1, "C", false, 0, "")
	pdf.Ln(20)
	pdf.SetXY(130, pdf.GetY())
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(60, 5, "(___________________)", "", 1, "C", false, 0, "")

	// Output buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper to format currency
func formatRupiah(amount int64) string {
	// Simple formatter
	// For production, use language/message package or similar
	s := fmt.Sprintf("%d", amount)
	// Insert dots
	n := len(s)
	if n <= 3 {
		return "Rp " + s
	}
	var result []byte
	for i, c := range s {
		if (n-i)%3 == 0 && i > 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return "Rp " + string(result)
}

func formatPeriod(period string) string {
	// 2026-01 -> Januari 2026
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return period
	}
	return t.Format("January 2006")
}
