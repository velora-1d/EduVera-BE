package export

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type ExportDomain interface {
	ExportStudents(ctx context.Context, tenantID string, format string) ([]byte, string, error)
	ExportPayments(ctx context.Context, tenantID string, format string) ([]byte, string, error)
}

type exportDomain struct {
	db outbound_port.DatabasePort
}

func NewExportDomain(db outbound_port.DatabasePort) ExportDomain {
	return &exportDomain{db: db}
}

// ExportStudents exports student list as PDF or Excel
func (d *exportDomain) ExportStudents(ctx context.Context, tenantID string, format string) ([]byte, string, error) {
	// Get student data
	students, err := d.db.Sekolah().GetSiswaByTenant(tenantID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get students: %w", err)
	}

	if format == "xlsx" {
		return d.exportStudentsExcel(students)
	}
	return d.exportStudentsPDF(students)
}

// ExportPayments exports payment report as PDF or Excel
func (d *exportDomain) ExportPayments(ctx context.Context, tenantID string, format string) ([]byte, string, error) {
	// Get payment/SPP data - using tabungan as example
	tabungan, err := d.db.Sekolah().GetTabunganList(tenantID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get payments: %w", err)
	}

	if format == "xlsx" {
		return d.exportPaymentsExcel(tabungan)
	}
	return d.exportPaymentsPDF(tabungan)
}

// ========== PDF Exports ==========

func (d *exportDomain) exportStudentsPDF(students []model.Siswa) ([]byte, string, error) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Arial", "", 10)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Laporan Data Siswa")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Tanggal: %s", time.Now().Format("02 Jan 2006")))
	pdf.Ln(10)

	// Table Header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(10, 8, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "NIS", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Nama", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Kelas", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 8, "Alamat", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Status", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 9)
	for i, s := range students {
		pdf.CellFormat(10, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 7, s.NIS, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 7, s.Nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, s.KelasNama, "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 7, s.Alamat, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, s.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 6, fmt.Sprintf("Total: %d siswa", len(students)))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("laporan_siswa_%s.pdf", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

func (d *exportDomain) exportPaymentsPDF(tabungan []model.Tabungan) ([]byte, string, error) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetFont("Arial", "", 10)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Laporan Keuangan - Tabungan Siswa")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Tanggal: %s", time.Now().Format("02 Jan 2006")))
	pdf.Ln(10)

	// Table Header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(10, 8, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 8, "Nama Santri", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Saldo", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Status", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Terakhir Update", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 9)
	var totalSaldo int64
	for i, t := range tabungan {
		pdf.CellFormat(10, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 7, t.NamaSantri, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("Rp %d", t.Saldo), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, t.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, t.UpdatedAt.Format("02 Jan 2006"), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
		totalSaldo += t.Saldo
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Total Saldo: Rp %d", totalSaldo))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("laporan_keuangan_%s.pdf", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

// ========== Excel Exports ==========

func (d *exportDomain) exportStudentsExcel(students []model.Siswa) ([]byte, string, error) {
	f := excelize.NewFile()
	sheet := "Siswa"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Set headers
	headers := []string{"No", "NIS", "Nama", "Kelas", "Alamat", "Status"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 5)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 30)
	f.SetColWidth(sheet, "D", "D", 20)
	f.SetColWidth(sheet, "E", "E", 30)
	f.SetColWidth(sheet, "F", "F", 12)

	// Data rows
	for i, s := range students {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), s.NIS)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.Nama)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.KelasNama)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.Alamat)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), s.Status)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("laporan_siswa_%s.xlsx", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

func (d *exportDomain) exportPaymentsExcel(tabungan []model.Tabungan) ([]byte, string, error) {
	f := excelize.NewFile()
	sheet := "Keuangan"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Set headers
	headers := []string{"No", "Nama Santri", "Saldo", "Status", "Terakhir Update"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 5)
	f.SetColWidth(sheet, "B", "B", 30)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 12)
	f.SetColWidth(sheet, "E", "E", 20)

	// Data rows
	var totalSaldo int64
	for i, t := range tabungan {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.NamaSantri)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.Saldo)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), t.UpdatedAt.Format("02 Jan 2006"))
		totalSaldo += t.Saldo
	}

	// Total row
	totalRow := len(tabungan) + 2
	f.SetCellValue(sheet, fmt.Sprintf("B%d", totalRow), "Total")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", totalRow), totalSaldo)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("laporan_keuangan_%s.xlsx", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}
