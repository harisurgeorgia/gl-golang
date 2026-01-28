// helpers.go (or in same controllers file)
package utils

import (
	"fmt"
	"gl/models"
	"gl/session"
	"gl/views"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBasePageData(c *gin.Context, title, header, page string, script *string) (views.PageData, error) {
	idStr := session.GetSession(c, "user_id")

	if idStr == "" {
		return views.PageData{
			Title:  title,
			Header: header,
			Search: true,
			User:   views.UserData{},
			Menus:  nil,
			Script: script,
			Link:   "",
		}, nil

	}
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

	menus, err := models.GetUserMenu(page)
	if err != nil {

	}

	/* 	menus := []views.UserMenu{
		{MenuDescription: "New", Icon: "fa fa-home", Url: "/journal", UserType: "editor", ItemType: "link", Page: "journal-entry"},
		{MenuDescription: "Edit Journal", Icon: "fa fa-users", Url: "/journal/list"},
		{MenuDescription: "Post Journal", Icon: "fa fa-book", Url: "/journal/list"},
		{MenuDescription: "Close Period", Icon: "fa fa-calendar", Url: "/close-period"},
	} */

	return views.PageData{
		Title:  title,
		Header: header,
		Search: true,
		User:   user,
		Menus:  menus,
		Script: script,
	}, nil
}
