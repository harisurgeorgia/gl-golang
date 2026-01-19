package controllers

import (
	"gl/models"
	"gl/reports"

	"github.com/gin-gonic/gin"
)

func TrialBalance(c *gin.Context) {
	p := 1
	rows, err := models.GetTrialBalance(p)
	if err != nil {
		c.HTML(500, "error.templ", nil)
		return
	}

	reports.GenerateTrialBalancePDF(rows, c)
	if err != nil {
		c.HTML(500, "error.templ", nil)
		return
	}

}
