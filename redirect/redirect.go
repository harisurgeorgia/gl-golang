package redirect

import (
	"gl/messages"

	"github.com/gin-gonic/gin"
)

var Error404 = messages.Error404
var Error403 = messages.Error403
var Error400 = messages.Error400
var Error500 = messages.Error500

func Init(c *gin.Context) {
	Error404.Redirect = c.Request.Referer()
	Error403.Redirect = c.Request.Referer()
	Error400.Redirect = c.Request.Referer()
	Error500.Redirect = c.Request.Referer()
}
