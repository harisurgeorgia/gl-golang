package controllers

import (
	"gl/messages"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func PageNotFound(c *gin.Context) {
	//var err error
	var data, err = utils.GetBasePageData(c, "GL/404", "", "", nil)
	if err != nil {

	}
	//var data = views.PageData{Title: "GL/404", Header: ""}
	utils.Render(c, 200, views.Layout(views.Nav(data), data, views.ErrorPage(messages.Error404)))
}
