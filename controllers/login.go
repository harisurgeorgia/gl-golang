package controllers

import (
	"database/sql"
	"fmt"
	"gl/db"
	"gl/mail"
	"gl/models"
	"gl/session"
	"gl/utils"
	"gl/validation"
	"gl/views"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Login(c *gin.Context) {

	var data = getLoginPageData()
	utils.Render(c, 200, views.Layout(data, views.LoginForm(data.Header, "")))
}

func LoginSubmit(c *gin.Context) {

	var data = getLoginPageData()
	// 1. Capture form input
	email := utils.NormalizeEmail(c.PostForm("email"))
	password := strings.TrimSpace(c.PostForm("password"))

	query := `SELECT id, password,  COALESCE(role, '') AS role, fullname FROM general_ledger.users WHERE email = $1 and is_active = true`
	row := db.Conn.QueryRow(query, email)

	//var id, hash, role, fullname string
	var user models.User

	err := row.Scan(&user.Id, &user.Password, &user.Role, &user.Fullname)
	if err != nil {
		log.Println("Error querying user:", err)
		utils.Render(c, http.StatusUnauthorized, views.Layout(
			data,
			views.LoginForm(data.Header, "Invalid credentials."),
		))
		return
	}
	if !utils.CheckPasswordHash(password, user.Password) {

		log.Println("Invalid password for user:", email)
		utils.Render(c, http.StatusUnauthorized, views.Layout(
			data,
			views.LoginForm(data.Header, "Invalid credentials."),
		))
		return
	}
	idStr := strconv.FormatInt(*user.Id, 10)
	session.SetSession(c, "user_id", idStr)
	session.SetSession(c, "user_email", user.Email)
	session.SetSession(c, "user_name", user.Fullname)
	session.SetSession(c, "user_role", *user.Role) // Assuming role is admin for simplicity

	c.Redirect(http.StatusFound, "/dashboard")
}

func getLoginPageData() views.PageData {
	return views.PageData{Title: "GL/Login", Header: "Login Page"}
}

func ForgotPassword(c *gin.Context) {
	var email string
	if c.Request.Method == http.MethodPost {
		email = utils.NormalizeEmail(c.PostForm("email"))
		err := validation.EmailValid(email)
		if err != nil {
			data := views.PageData{Title: "GL/Forgot Password", Header: "Reset Password"}
			utils.Render(c, 200, views.Layout(data, views.ResetFrom(data.Header, err.Error(), email)))
			return
		} else {
			token, err := utils.GenerateResetToken()
			if err != nil {
				c.String(http.StatusInternalServerError,
					"Could not generate reset token: %v", err,
				)
				return
			}

			link := fmt.Sprintf(`http://localhost:8080/change-password/%s`, token)
			body := fmt.Sprintf(`<p>Click <a href="%s">here</a> to reset your password.</p>`, link)
			if err := mail.SendMail(email, "Password Reset", body); err != nil {
				c.String(http.StatusInternalServerError, "Could not send email: %v", err)
				return
			}
			err = saveResetToken(email, token)
			if err != nil {
				c.String(http.StatusInternalServerError, "Something went wrong %s", err)
			}
			c.String(http.StatusOK, "Reset link sent!")

			return
		}
	} else {
		data := views.PageData{Title: "GL/Forgot Password", Header: "Reset Password"}
		utils.Render(c, 200, views.Layout(data, views.ResetFrom(data.Header, "", "")))
	}

}

func ChangePassword(c *gin.Context) {
	data := views.PageData{Title: "GL/Reset Password", Header: "Reset Password"}
	if c.Request.Method == http.MethodGet {
		token := strings.TrimSpace(strings.TrimPrefix(c.Param("key"), "/"))
		utils.Render(c, 200, views.Layout(data, views.ChangePasswordForm(data.Header, "", "", token, "", "")))
		return
	}
	email := utils.NormalizeEmail(c.PostForm("email"))
	err := validation.EmailValid(email)

	if err != nil {
		utils.Render(c, http.StatusSeeOther, views.Layout(data, views.ChangePasswordForm(data.Header, err.Error(), email, c.PostForm("token"), c.PostForm("password"), c.PostForm("confirm-password"))))
	}

	err = validation.CheckPasswordMatch(c.PostForm("password"), c.PostForm("confirm-password"))

	if err != nil {
		utils.Render(c, http.StatusSeeOther, views.Layout(data, views.ChangePasswordForm(data.Header, err.Error(), email, c.PostForm("token"), c.PostForm("password"), c.PostForm("confirm-password"))))
	}
	//var hashPassword string

	hashPassword, err := utils.HashPassword(strings.TrimSpace(c.PostForm("password")))
	if err != nil {
		utils.Render(c, http.StatusSeeOther, views.Layout(data, views.ChangePasswordForm(data.Header, "unknown token or email", email, c.PostForm("token"), c.PostForm("password"), c.PostForm("confirm-password"))))
		return
	}
	var result sql.Result
	result, err = db.Conn.Exec(`
        UPDATE general_ledger.users
           SET reset_token = $1,
		   password = $2
         WHERE email  = $3 
		 and reset_token = $4    
    `, "", hashPassword, email, c.PostForm("token"))

	if err != nil {
		utils.Render(c, http.StatusSeeOther, views.Layout(data, views.ChangePasswordForm(data.Header, "unknown token or email", email, c.PostForm("token"), c.PostForm("password"), c.PostForm("confirm-password"))))
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		utils.Render(c, http.StatusSeeOther, views.Layout(data, views.ChangePasswordForm(data.Header, "unknown token or email", email, c.PostForm("token"), c.PostForm("password"), c.PostForm("confirm-password"))))
	}
	c.Redirect(http.StatusOK, "/")

}

func saveResetToken(email, token string) error {
	// Use Exec because we don’t need to return any rows
	result, err := db.Conn.Exec(`
        UPDATE general_ledger.users
           SET reset_token = $1
         WHERE email       = $2
    `, token, email)
	if err != nil {
		return fmt.Errorf("no user found%w", err)
	}
	// Optionally, check that a row was actually updated:
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no user found with email %q", email)
	}
	return nil
}

func Logout(c *gin.Context) {
	session.LogoutHandler(c)
	c.Redirect(http.StatusSeeOther, "/")
}

// utils.Render is a utility function to render templates

// shared render function
