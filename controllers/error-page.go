package controllers

import (
	"gl/messages"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func PageNotFound(c *gin.Context) {
	//var err error
	var data, err = getBasePageData(c, "GL/404", "", "")
	if err != nil {

	}
	//var data = views.PageData{Title: "GL/404", Header: ""}
	utils.Render(c, 200, views.Layout(views.Nav(nil, false), data, views.ErrorPage(messages.Error404)))
}
