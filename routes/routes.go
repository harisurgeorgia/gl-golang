package routes

import (
	"gl/controllers"
	"gl/middleware"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// Group with middleware
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/user", controllers.UserCreate)
		authGroup.POST("/user", controllers.UserSave)
		authGroup.GET("/user/:id", controllers.GetUser)
		authGroup.GET("/dashboard", controllers.Dashboard)
		authGroup.GET("/journal", controllers.JournalEntry)
		authGroup.GET("/journal/list", controllers.JournalList)
		authGroup.POST("/journal/save", controllers.JournalSave)
		authGroup.GET("/close-period", controllers.ClosePeriod)
		authGroup.GET("/logout", controllers.Logout)
	}
	// Public routes
	r.GET("/", middleware.RedirectIfAuthenticated(), controllers.Login)
	r.POST("/", controllers.LoginSubmit)
	// route for page not found
	r.NoRoute(func(c *gin.Context) {
		utils.Render(c, 404, views.Layout(views.PageData{
			Title:  "Page Not Found",
			Header: "404 - Page Not Found",
		}, views.View404()))
	})

	// route for unexpected error
	r.GET("/unexpected-error", func(c *gin.Context) {
		c.HTML(500, "500.templ", nil)
	})

}
