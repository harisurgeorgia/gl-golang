package controllers

import (
	"fmt"
	"gl/models"
	"gl/redirect"
	"gl/utils"
	"gl/views"

	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func JournalSearch(c *gin.Context) {
	search := strings.TrimSpace(c.PostForm("search"))

	if search != "" {

	}

	id, err := models.FindJournalByJournalNumber(search)
	if err != nil {

	}
	fmt.Println("/journal/edit/" + strconv.FormatInt(id, 10))
	c.Redirect(http.StatusFound, "/journal/edit/"+strconv.FormatInt(id, 10))

}

func FindProductByCode(c *gin.Context) {

	search, exists := c.GetPostForm("search")
	if !exists || strings.TrimSpace(search) == "" {
		utils.Render(c, 400, views.Layout(nil, views.PageData{
			Title: "Bad Request",
		}, views.ErrorPage(redirect.Error400)))
		return
	}

	products, err := models.FindProductByCode(search)
	if err != nil || len(products) == 0 {
		utils.Render(c, 404, views.Layout(nil, views.PageData{
			Title: "Product Not Found",
		}, views.ErrorPage(redirect.Error404)))
		return
	}
	if len(products) == 1 {
		fmt.Println("/inventory/product/" + strconv.FormatInt(*products[0].Id, 10))
		c.Redirect(http.StatusFound, "/inventory/product/"+strconv.FormatInt(*products[0].Id, 10))
	} else {
		c.Redirect(http.StatusFound, "/inventory/product/list-products/?filter="+search)
	}
	//

}
