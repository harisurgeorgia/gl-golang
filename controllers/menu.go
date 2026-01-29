package controllers

import (
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func Menu(header, menus string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := utils.GetBasePageData(c, header, header, menus, nil)
		data.Search = false

		utils.Render(
			c,
			200,
			views.Layout(
				views.Nav(data),
				data,
				views.MenuPage(data),
			),
		)
	}
}
