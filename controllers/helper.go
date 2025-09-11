// helpers.go (or in same controllers file)
package controllers

import (
	"fmt"
	"gl/session"
	"gl/views"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getBasePageData(c *gin.Context, title, header string) (views.PageData, error) {
	idStr := session.GetSession(c, "user_id")
	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return views.PageData{}, fmt.Errorf("invalid user_id in session: %w", err)
	}

	user := views.UserData{
		Id:       idInt64,
		Email:    session.GetSession(c, "user_email"),
		Fullname: session.GetSession(c, "user_name"),
		Role:     session.GetSession(c, "user_role"),
	}

	menus := []views.UserMenu{
		{MenuDescription: "Create Journal", Icon: "fa fa-home", Url: "/journal", UserType: "editor", ItemType: "button", Page: "journal-entry"},
		{MenuDescription: "Edit Journal", Icon: "fa fa-users", Url: "/journal/list"},
		{MenuDescription: "Post Journal", Icon: "fa fa-book", Url: "/journal/list"},
		{MenuDescription: "Close Period", Icon: "fa fa-calendar", Url: "/close-period"},
	}

	return views.PageData{
		Title:  title,
		Header: header,
		User:   user,
		Menus:  menus,
		Script: "/static/js/journal-entry.js",
	}, nil
}
