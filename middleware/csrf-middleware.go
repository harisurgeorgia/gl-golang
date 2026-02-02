package middleware

import (
	"gl/redirect"
	"gl/utils"
	"gl/views"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func CSRFMiddleware(r *gin.Engine) {
	r.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			utils.Render(c, 400, views.Layout(nil, views.LayoutAttribute{
				PageTitle: "CSRF Token Mismatch",
			}, views.ErrorPage(redirect.Error400)))
			c.Abort()
		},
	}))
}
