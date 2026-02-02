package routes

import (
	"gl/Inventory"
	"gl/controllers"
	"gl/messages"
	"gl/middleware"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	middleware.CSRFMiddleware(r)
	//r.GET("/user", controllers.UserCreate)
	//r.POST("/user", controllers.UserSave)
	// Group with middleware
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/user", controllers.UserCreate)
		authGroup.POST("/user", controllers.UserSave)
		authGroup.GET("/user/:id", middleware.RequireRole(3), controllers.GetUser)
		authGroup.GET("/dashboard", controllers.Menu("GL/Dashboard", "dashboard"))
		authGroup.GET("/user-menu", controllers.Menu("GL/User Menu", "user"))
		authGroup.GET("/journal-menu", controllers.Menu("GL/Journal Menu", "journal"))
		authGroup.GET("/product-menu", controllers.Menu("GL/Product Menu", "product"))
		authGroup.GET("/journal", middleware.RequireRole(3), controllers.JournalEntry)
		authGroup.GET("/journal/list-by-status/:filter", controllers.JournalList)
		authGroup.GET("/journal/list-by-number/:filter", controllers.FindJournalByJournalNumber)
		authGroup.GET("/journal/edit/:id", controllers.JournalEdit)
		authGroup.POST("/journal/save", controllers.JournalSave)
		authGroup.POST("/journal/post", controllers.JournalPost)
		authGroup.GET("/close-period", controllers.ClosePeriod)
		authGroup.GET("/logout", controllers.Logout)
		authGroup.GET("/journal/search", controllers.JournalSearch)
		authGroup.GET("/period-list", controllers.PeriodList)
		authGroup.GET("/period/:id", controllers.PeriodAddEdit)
		authGroup.GET("/period", controllers.PeriodAddEdit)
		authGroup.POST("/period", controllers.PeriodSave)
		authGroup.GET("verify", controllers.JournalVerify)
		authGroup.GET("/trial-balance", controllers.TrialBalance)
		authGroup.GET("/general-ledger", controllers.GeneralLedger)
		authGroup.GET("/inventory/product", Inventory.Product)
		authGroup.GET("/inventory/product/:id", Inventory.Product)
		authGroup.POST("/inventory/product", Inventory.SaveProduct)
		authGroup.GET("/inventory/product/search", controllers.FindProductByCode)
		authGroup.GET("/inventory/product/list-products", Inventory.ListAllProducts)
		authGroup.GET("/inventory/product/list-products/:filter", Inventory.ListAllProducts)
		authGroup.DELETE("/inventory/product/:id", Inventory.DeleteProduct)
		authGroup.GET("/csrf", controllers.FormCSRF)
		authGroup.POST("/csrf", controllers.FormCSRF)

	}
	// Public routesp
	r.GET("/", middleware.RedirectIfAuthenticated(), controllers.Login)
	r.POST("/", controllers.LoginSubmit)
	// route for page not found
	r.NoRoute(func(c *gin.Context) {
		utils.Render(c, 404, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Page Not Found",
		}, views.ErrorPage(messages.Error404)))
	})

	// route for unexpected error
	r.GET("/unexpected-error", func(c *gin.Context) {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
	})

}
