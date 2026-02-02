package Inventory

import (
	"gl/messages"
	"gl/models"
	"gl/redirect"
	"gl/session"
	"gl/utils"
	"gl/views"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func Product(c *gin.Context) {
	product := &models.Product{}
	idStr := c.Param("id")
	trimmed := strings.TrimSpace(idStr)
	s := "/inventory/product/search"
	pageData := views.LayoutAttribute{
		PageTitle: "Product",
		Script:    nil,
	}
	menus, err := models.GetUserMenu("dashboard")
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}
	navData := views.NavData{
		Link:  &s,
		Menus: menus,
	}

	grouped, err := models.GetAllSubAccountType()
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}
	token := csrf.GetToken(c)

	if len(trimmed) == 0 {
		utils.Render(c, http.StatusOK, views.Layout(views.Nav(navData), pageData, views.Product(product, grouped, models.Units, token)))
		return
	}

	id, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

	p, err := models.GetProduct(id)

	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(messages.Error500)))
		return
	}

	utils.Render(c, 200, views.Layout(views.Nav(navData), pageData, views.Product(p, grouped, models.Units, token)))
}

func SaveProduct(c *gin.Context) {
	product := &models.Product{}
	val := session.GetSession(c, "user_id")

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

	c.ShouldBind(product)
	product.UserId = userID
	id, err := models.SaveProduct(*product)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}
	c.Redirect(http.StatusSeeOther, "/inventory/product/"+strconv.FormatInt(*id, 10))
}

func UpdateProduct(c *gin.Context) {

}

func DeleteProduct(c *gin.Context) {

}

func ListAllProducts(c *gin.Context) {

	menus, err := models.GetUserMenu("dashboard")
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}

	link := "/inventory/product/search"
	navData := views.NavData{
		Link:   &link,
		Menus:  menus,
		Search: true,
	}

	filter := c.Param("filter")
	if len(filter) == 0 {
		filter = ""
	}
	product, err := models.ListAllProducts(filter)
	if err != nil {
		utils.Render(c, 500, views.Layout(nil, views.LayoutAttribute{
			PageTitle: "Unexpected Error",
		}, views.ErrorPage(redirect.Error500)))
		return
	}
	utils.Render(c, 200, views.Layout(views.Nav(navData), views.LayoutAttribute{
		PageTitle: "Product List",
		Script:    nil,
	}, views.ProductList(product)))

}
