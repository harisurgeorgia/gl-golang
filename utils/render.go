package utils

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, status int, component templ.Component) {
	c.Status(status)
	_ = component.Render(c.Request.Context(), c.Writer)
}
