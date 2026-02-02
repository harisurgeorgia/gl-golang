package Inventory

import (
	"gl/models"
	"gl/redirect"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func ProductPurchase(c *gin.Context) {
	link := "/inventory/product/search"
	menus, err := models.GetUserMenu("dashboard")
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

	navData := views.NavData{
		Link:   &link,
		Menus:  menus,
		Search: true,
	}
	invoice := models.Invoice{}
	token := csrf.GetToken(c)

	utils.Render(
		c,
		200,
		views.Layout(
			views.Nav(navData),
			views.LayoutAttribute{
				PageTitle: "Product Purchase",
				Script:    nil,
			},
			views.ProductForm(invoice, token),
		),
	)
}

func SaveProductPurchase(c *gin.Context) {

}

func ProductPurchaseEdit(c *gin.Context) {

}
