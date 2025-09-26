// views/types.go
package views

import "gl/models"

type PageData struct {
	Title  string
	Header string
	Role   string
	Search bool
	User   UserData
	Menus  []models.UserMenu
	Script string
}
type UserData struct {
	Id       int64
	Email    string
	Fullname string
	Password string
	Role     string
}
