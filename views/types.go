// views/types.go
package views

import "gl/models"

type PageData struct {
	User        *UserData
	RequestLink *string
}

type UserData struct {
	Id       int64
	Email    string
	Fullname string
	Password string
	Role     string
}

type LayoutAttribute struct {
	PageTitle  string
	PageHeader string
	Script     *string
}

type NavData struct {
	Menus     []models.UserMenu
	Search    bool
	Link      *string
	IsTopNav  bool
	CsrfToken string
}
