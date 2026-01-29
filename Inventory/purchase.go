package Inventory

import (
	"gl/models"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
)

func ProductPurchase(c *gin.Context) {
	data, _ := utils.GetBasePageData(c, "Product Purchase", "Product Purchase", "product-purchase", nil)
	data.Search = false
	invoice := models.Invoice{}

	utils.Render(
		c,
		200,
		views.Layout(
			views.Nav(data),
			data,
			views.ProductForm(invoice),
		),
	)
}

func SaveProductPurchase(c *gin.Context) {

}

func ProductPurchaseEdit(c *gin.Context) {

}
