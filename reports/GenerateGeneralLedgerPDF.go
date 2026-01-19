package reports

import (
	"gl/models"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func GenerateGeneralLedgerPDF(rows []models.GeneralLedgerRow, closingBalance *models.ClosingBalance, c *gin.Context) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "General Ledger")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	for _, row := range rows {
		pdf.Cell(0, 10, row.JournalNumber)
		pdf.Ln(10)
	}
}
