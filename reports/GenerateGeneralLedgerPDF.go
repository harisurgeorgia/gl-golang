package reports

import (
	"bytes"
	"gl/models"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GenerateGeneralLedgerPDF(rows []models.GeneralLedgerRow, closingBalance *models.ClosingBalance, c *gin.Context) error {

	const bottomLimit = 270.0
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	openingBalance(pdf, closingBalance)
	for _, row := range rows {
		pdf.CellFormat(40, 7, row.JournalNumber, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 7, row.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, FormatCurrency(row.Debit, language.AmericanEnglish, currency.USD), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, FormatCurrency(row.Credit, language.AmericanEnglish, currency.USD), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, FormatCurrency(row.Balance, language.AmericanEnglish, currency.USD), "1", 1, "R", false, 0, "")
		if pdf.GetY()+7 > bottomLimit {
			generateGLHeader(pdf)
		}
	}
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

func openingBalance(pdf *gofpdf.Fpdf, closingBalance *models.ClosingBalance) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "General Ledger", "1", 1, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(20, 7, closingBalance.AccountCode, "1", 0, "L", true, 0, "")
	pdf.CellFormat(55, 7, closingBalance.AccountName, "1", 0, "L", true, 0, "")
	pdf.CellFormat(55, 7, closingBalance.Description, "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, FormatCurrency(closingBalance.Debit, language.AmericanEnglish, currency.USD), "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, FormatCurrency(closingBalance.Credit, language.AmericanEnglish, currency.USD), "1", 1, "R", true, 0, "")

}

func FormatCurrency(d float64, tag language.Tag, curr currency.Unit) string {
	// Example locale: en-US

	p := message.NewPrinter(tag)

	return p.Sprintf("%s%.2f", currency.Symbol(curr), d)
}

func generateGLHeader(pdf *gofpdf.Fpdf) {
	pdf.Ln(10)

	pdf.CellFormat(40, 7, "Journal Number", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 7, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Debit", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, "Credit", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, "Balance", "1", 1, "R", true, 0, "")
}
