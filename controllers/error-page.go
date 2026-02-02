package controllers

import (
	"gl/messages"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func PageNotFound(c *gin.Context) {
	//var err error
	attr := views.LayoutAttribute{
		PageTitle:  "GL/404",
		PageHeader: "GL/404",
		Script:     nil,
	}

	utils.Render(c, 200, views.Layout(nil, attr, views.ErrorPage(messages.Error404)))
}
