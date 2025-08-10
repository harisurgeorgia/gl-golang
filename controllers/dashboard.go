package controllers

import (
	"gl/session"
	"gl/utils"
	"gl/views"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {

	idStr := session.GetSession(c, "id")
	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid session id: %v", idStr)
		c.Redirect(http.StatusBadRequest, "/unexpected-error")
		return
	}

	user := views.UserData{
		Id:       idInt64,
		Email:    session.GetSession(c, "user_email"),
		Fullname: session.GetSession(c, "user_name"),
		Role:     session.GetSession(c, "user_role"),
		Password: "", // probably don't store password in struct from session
	}

	data := views.PageData{
		Title:  "Dashboard",
		Header: "Dashboard",
		User:   user,
	}

	utils.Render(c, 200, views.Layout(data, views.DashboardPage(data))) // Assuming msg is a string variable with a welcome message
}
