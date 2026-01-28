package Inventory

import (
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func ProductPurchase(c *gin.Context) {
	data, _ := utils.GetBasePageData(c, "Product Purchase", "Product Purchase", "product-purchase", nil)
	data.Search = false

	utils.Render(
		c,
		200,
		views.Layout(
			views.Nav(data.Menus, data.Search, &data.Link),
			data,
			views.ProductForm(data),
		),
	)
}

func SaveProductPurchase(c *gin.Context) {

}

func ProductPurchaseEdit(c *gin.Context) {

}
