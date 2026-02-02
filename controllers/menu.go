package controllers

import (
	"gl/models"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func Menu(header, menu string) gin.HandlerFunc {
	return func(c *gin.Context) {

		menus, _ := models.GetUserMenu(menu)
		navData := views.NavData{
			Menus:    menus,
			Link:     nil,
			IsTopNav: false,
		}

		utils.Render(
			c,
			200,
			views.Layout(
				views.Nav(navData),
				views.LayoutAttribute{PageTitle: header, Script: nil},
				views.MenuPage(navData, header),
			),
		)
	}
}
