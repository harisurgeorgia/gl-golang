package controllers

import (
	"gl/messages"
	"gl/models"
	"gl/session"
	"gl/utils"
	"gl/views"

	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Product(c *gin.Context) {
	product := &models.Product{}
	idStr := c.Param("id")
	trimmed := strings.TrimSpace(idStr)
	data, err := getBasePageData(c, "Product", "Product Edit/Entry", "", nil)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}

	grouped, err := models.GetAllSubAccountType()
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}

	if len(trimmed) == 0 {
		utils.Render(c, http.StatusOK, views.Layout(views.Nav(data.Menus, data.Search), data, views.Product(product, grouped, models.Units)))
		return
	}

	id, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}

	p, err := models.GetProduct(id)

	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}
	utils.Render(c, 200, views.Layout(views.Nav(data.Menus, data.Search), data, views.Product(p, grouped, models.Units)))
}

func SaveProduct(c *gin.Context) {
	product := &models.Product{}
	val := session.GetSession(c, "user_id")

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}

	c.ShouldBind(product)
	product.UserId = userID
	id, err := models.SaveProduct(*product)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.PageData{
			Title: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}
	c.Redirect(http.StatusSeeOther, "/inventory/product/"+strconv.FormatInt(*id, 10))
}

func UpdateProduct(c *gin.Context) {

}

func DeleteProduct(c *gin.Context) {

}

func ListAllProducts(c *gin.Context) {

}
