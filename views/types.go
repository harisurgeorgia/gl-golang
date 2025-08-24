// views/types.go
package views

type PageData struct {
	Title  string
	Header string
	Role   string
	User   UserData
	Menus  []UserMenu
}
type UserData struct {
	Id       int64
	Email    string
	Fullname string
	Password string
	Role     string
}

type UserMenu struct {
	Id              int64
	MenuDescription string
	Url             string
	Icon            string
	Active          bool
	IsAdmin         bool
}
