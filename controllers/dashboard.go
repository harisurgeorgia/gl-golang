package controllers

import (
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {

	data, _ := getBasePageData(c, "Dashboard", "Dashboard", "dashboard")
	data.Search = false
	utils.Render(c, 200, views.Layout(views.Nav(data.Menus, data.Search), data, views.DashboardPage(data))) // Assuming msg is a string variable with a welcome message
}
