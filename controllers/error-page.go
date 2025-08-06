package controllers

import (
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func PageNotFound(c *gin.Context) {
	var data = views.PageData{Title: "GL/404", Header: ""}
	utils.Render(c, 200, views.Layout(data, views.View404()))
}
