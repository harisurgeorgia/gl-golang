package routes

import (
	"gl/controllers"
	"gl/messages"
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
		authGroup.GET("/user/:id", middleware.RequireRole(3), controllers.GetUser)
		authGroup.GET("/dashboard", controllers.Menu("My Dashboard", "dashboard"))
		authGroup.GET("/user-menu", controllers.Menu("User Menu", "user"))
		authGroup.GET("/journal-menu", controllers.Menu("Journal Menu", "journal"))
		authGroup.GET("/journal", middleware.RequireRole(3), controllers.JournalEntry)
		authGroup.GET("/journal/list/:filter", controllers.JournalList)
		authGroup.GET("/journal/edit/:id", controllers.JournalEdit)
		authGroup.POST("/journal/save", controllers.JournalSave)
		authGroup.POST("/journal/post", controllers.JournalPost)
		authGroup.GET("/close-period", controllers.ClosePeriod)
		authGroup.GET("/logout", controllers.Logout)
		authGroup.POST("/search", controllers.Search)
		authGroup.GET("/period-list", controllers.PeriodList)
		authGroup.GET("/period/:id", controllers.PeriodAddEdit)
		authGroup.GET("/period", controllers.PeriodAddEdit)
		authGroup.POST("/period", controllers.PeriodSave)
		authGroup.GET("verify", controllers.JournalVerify)
		authGroup.GET("/trial-balance", controllers.TrialBalance)

	}
	// Public routes
	r.GET("/", middleware.RedirectIfAuthenticated(), controllers.Login)
	r.POST("/", controllers.LoginSubmit)
	// route for page not found
	r.NoRoute(func(c *gin.Context) {
		utils.Render(c, 404, views.Layout(nil, views.PageData{
			Title: "Page Not Found",
		}, views.ErrorPage(messages.Error404)))
	})

	// route for unexpected error
	r.GET("/unexpected-error", func(c *gin.Context) {
		c.HTML(500, "500.templ", nil)
	})

}
