package controllers

import (
	"fmt"
	"gl/db"
	"gl/models"
	"gl/redirect"
	"gl/utils"
	"gl/validation"
	"gl/views"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var user models.User

func UserCreatePage(c *gin.Context) {

	link := "/user/search"
	menus, _ := models.GetUserMenu("dashboard")
	navData := views.NavData{
		Link:   &link,
		Menus:  menus,
		Search: true,
	}
	utils.Render(c, 200, views.Layout(views.Nav(navData), views.LayoutAttribute{
		PageTitle: "User Information",
		Script:    nil,
	}, views.UserForm(user)))
}
func UserCreate(c *gin.Context) {

	link := "/user/search"
	menus, _ := models.GetUserMenu("dashboard")
	navData := views.NavData{
		Link:   &link,
		Menus:  menus,
		Search: true,
	}

	utils.Render(c, 200, views.Layout(views.Nav(navData), views.LayoutAttribute{
		PageTitle: "User Information",
		Script:    nil,
	}, views.UserForm(user)))
}
func UserSave(c *gin.Context) {

	email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	if (c.PostForm("id") != "") && (c.PostForm("id") != "0") {

		_, err := db.Conn.Exec(`
		UPDATE general_ledger.users
		SET email = $1, fullname = $2, updated_at = $3
		WHERE id = $4`, email, c.PostForm("fullname"), time.Now(), c.PostForm("id"))

		if err != nil {
			log.Println("Update error:", err)
			c.String(500, "Database update failed")
			return
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/user/%s", c.PostForm("id")))
		return
	}
	menus, _ := models.GetUserMenu("dashboard")
	link := "/user/search"
	navData := views.NavData{
		Link:   &link,
		Menus:  menus,
		Search: true,
	}
	var user models.User

	// Automatically fills fields from POST form data
	if err := c.ShouldBind(&user); err != nil {
		utils.Render(c, 400, views.Layout(views.Nav(navData), views.LayoutAttribute{
			PageTitle: "User Information",
			Script:    nil,
		}, views.UserForm(user)))
		//c.String(http.StatusBadRequest, "Invalid form input: %v", err)
		return
	}

	err := validation.EmailValid(email)
	if err != nil {
		utils.Render(c, 400, views.Layout(views.Nav(navData), views.LayoutAttribute{
			PageTitle: "User Information",
			Script:    nil,
		}, views.UserForm(user)))
	}

	err = validation.IsValidPassword(user.Password)
	if err != nil {
		utils.Render(c, 400, views.Layout(views.Nav(navData), views.LayoutAttribute{
			PageTitle: "User Information",
			Script:    nil,
		}, views.UserForm(user)))
		return
	}

	err = validation.CheckPasswordMatch(c.PostForm("password"), c.PostForm("confirmPassword"))

	if err != nil {
		utils.Render(c, 400, views.Layout(views.Nav(navData), views.LayoutAttribute{
			PageTitle: "User Information",
			Script:    nil,
		}, views.UserForm(user)))
		return
	}

	hash, err := utils.HashPassword(strings.TrimSpace(c.PostForm("password")))
	user.Email = email
	//bytes, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hash failed:", err)
		return
	}
	user.Password = hash

	err = db.Conn.QueryRow(`
		INSERT INTO general_ledger.users (email, fullname, password)
		VALUES ($1, $2, $3) RETURNING id
	`, user.Email, user.Fullname, user.Password).Scan(&user.Id)

	if err != nil {
		log.Println("DB insert error:", err)
		utils.Render(c, 500, views.Layout(views.Nav(navData), views.LayoutAttribute{
			PageTitle: "User Information",
			Script:    nil,
		}, views.UserForm(user)))
		return
	}
	log.Println("User created with ID:", user.Id)

	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d", user.Id))

}
func GetUser(c *gin.Context) {

	id := c.Param("id")
	query := `SELECT id, email, fullname FROM general_ledger.users WHERE id = $1`
	errdb := db.Conn.QueryRow(query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Fullname,
	)

	if errdb != nil {
		log.Println("No user found with ID:", id)
	}
	menus, err := models.GetUserMenu("dashboard")
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}
	link := "/user/search"
	navData := views.NavData{
		Link:  &link,
		Menus: menus,
	}
	utils.Render(c, 200, views.Layout(views.Nav(navData), views.LayoutAttribute{
		PageTitle: "User Information",
		Script:    nil,
	}, views.UserForm(user)))
}
