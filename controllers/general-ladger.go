package controllers

import (
	"gl/models"
	"gl/redirect"
	"gl/reports"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func GeneralLedger(c *gin.Context) {
	/* period := c.Query("period")
	periodInt, err := strconv.Atoi(period)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	} */
	rows, closingBalance, err := models.GetGeneralLedger(1, 5)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

	err = reports.GenerateGeneralLedgerPDF(rows, closingBalance, c)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

}
