package reports

import (
	"bytes"
	"fmt"
	"gl/models"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func formatMoney(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func GenerateTrialBalancePDF(rows []models.TrialBalanceRow, c *gin.Context) error {
	const bottomLimit = 270.0
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	generateHeader(pdf)

	var totalDebit, totalCredit float64

	for _, r := range rows {
		pdf.CellFormat(30, 7, r.AccountCode, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 7, r.AccountName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, formatMoney(r.Debit), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 7, formatMoney(r.Credit), "1", 1, "R", false, 0, "")

		totalDebit += r.Debit
		totalCredit += r.Credit
		if pdf.GetY()+7 > bottomLimit {
			generateHeader(pdf)
		}
	}

	if pdf.GetY()+7 > bottomLimit {
		generateHeader(pdf)
	}
	// Totals
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(100, 7, "Total", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 7, formatMoney(totalDebit), "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 7, formatMoney(totalCredit), "1", 1, "R", true, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return err
	}

	// Send to browser
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=trial_balance.pdf")
	c.Data(200, "application/pdf", buf.Bytes())

	return nil
}

func generateHeader(pdf *gofpdf.Fpdf) {

	// Table header
	pdf.AddPage()

	// Title
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 8, "Trial Balance", "", 1, "C", false, 0, "")

	pdf.Ln(4)

	// Table header
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(30, 7, "Account", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 7, "Name", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 7, "Debit", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 7, "Credit", "1", 1, "R", true, 0, "")

	pdf.SetFont("Helvetica", "", 10)
}
