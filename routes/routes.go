package routes

import (
	"gl/controllers"
	"gl/session"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginSubmit)                      // Assuming you want to handle POST requests as well
	r.GET("/user", session.AuthRequired(), controllers.UserCreate) // Example route for user creation page
	r.POST("/user", session.AuthRequired(), controllers.UserSave)  // Example route for user creation submission
	r.GET("/user/:id", session.AuthRequired(), controllers.GetUser)
	r.GET("/logout", session.LogoutHandler)
	r.GET("/forgot-password", controllers.ForgotPassword)
	r.POST("forgot-password", controllers.ForgotPassword)
	r.GET("/change-password/:key", controllers.ChangePassword)
	r.POST("/change-password", controllers.ChangePassword)
	r.NoRoute(controllers.PageNotFound)
	r.GET("/journal", controllers.JournalEntry)
	r.POST("/journal/save", controllers.JournalSave)
	r.GET("/close-period", controllers.ClosePeriod)
}
