package controllers

import (
	"fmt"
	"gl/models"

	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	search := strings.TrimSpace(c.PostForm("search"))

	if search != "" {

	}

	id, err := models.FindJournalByJournalNumber(search)
	if err != nil {

	}
	fmt.Println("/journal/edit/" + strconv.FormatInt(id, 10))
	c.Redirect(http.StatusFound, "/journal/edit/"+strconv.FormatInt(id, 10))

}
