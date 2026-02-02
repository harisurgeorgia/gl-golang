package controllers

import (
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func FormCSRF(c *gin.Context) {
	attr := views.LayoutAttribute{
		PageTitle:  "CSRF",
		PageHeader: "CSRF",
		Script:     nil,
	}

	token := csrf.GetToken(c)
	utils.Render(c, 200, views.Layout(nil, attr, views.CSRFPage(token)))
}
