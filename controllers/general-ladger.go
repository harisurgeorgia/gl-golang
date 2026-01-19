package controllers

import (
	"gl/models"
	"gl/reports"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GeneralLedger(c *gin.Context) {
	period := c.Query("period")
	periodInt, err := strconv.Atoi(period)
	if err != nil {
		c.HTML(500, "error.templ", nil)
		return
	}
	rows, closingBalance, err := models.GetGeneralLedger(periodInt, 1)
	if err != nil {
		c.HTML(500, "error.templ", nil)
		return
	}

	reports.GenerateGeneralLedgerPDF(rows, closingBalance, c)
	if err != nil {
		c.HTML(500, "error.templ", nil)
		return
	}

}
