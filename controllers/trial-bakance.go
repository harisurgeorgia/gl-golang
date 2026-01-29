package controllers

import (
	"gl/models"
	"gl/redirect"
	"gl/reports"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func TrialBalance(c *gin.Context) {
	p := 1
	rows, err := models.GetTrialBalance(p)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

	reports.GenerateTrialBalancePDF(rows, c)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

}
