package controllers

import (
	"fmt"
	"gl/db"
	"gl/models"
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
	var data = getUserPageData()
	utils.Render(c, 200, views.Layout(data, views.UserForm(data.Header, "", user)))
}
func UserCreate(c *gin.Context) {
	var data = getUserPageData()
	utils.Render(c, 200, views.Layout(data, views.UserForm(data.Header, "", user)))
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
	var data = getUserPageData()
	var user models.User

	// Automatically fills fields from POST form data
	if err := c.ShouldBind(&user); err != nil {
		utils.Render(c, 400, views.Layout(data, views.UserForm(data.Header, "Required fields are missing or invalid.", user)))
		//c.String(http.StatusBadRequest, "Invalid form input: %v", err)
		return
	}

	err := validation.EmailValid(email)
	if err != nil {
		utils.Render(c, 400, views.Layout(data, views.UserForm(data.Header, err.Error(), user)))
	}

	err = validation.IsValidPassword(user.Password)
	if err != nil {
		utils.Render(c, 400, views.Layout(data, views.UserForm(data.Header, err.Error(), user)))
		return
	}

	err = validation.CheckPasswordMatch(c.PostForm("password"), c.PostForm("confirm-password"))

	if err != nil {
		utils.Render(c, 400, views.Layout(data, views.UserForm(data.Header, "Passwords do not match or are empty.", user)))
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
		utils.Render(c, 500, views.UserForm("Registration", "Could not save user", user))
		return
	}
	log.Println("User created with ID:", user.Id)

	c.Redirect(http.StatusFound, fmt.Sprintf("/user/%d", user.Id))

}
func GetUser(c *gin.Context) {

	var data = getUserPageData()
	id := c.Param("id")
	query := `SELECT id, email, fullname FROM general_ledger.users WHERE id = $1`
	err := db.Conn.QueryRow(query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Fullname,
	)

	if err != nil {
		log.Println("No user found with ID:", id)
	}

	utils.Render(c, 200, views.Layout(data, views.UserForm(data.Header, "", user)))
}

func getUserPageData() views.PageData {
	return views.PageData{Title: "GL/Maintenence", Header: "User Information"}
}
