package middleware

import (
	"gl/messages"
	"gl/utils"
	"gl/views"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireRole(minRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)

		userID := sess.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		userTypeAny := sess.Get("user_role")
		if userTypeAny == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		userTypeStr, ok := userTypeAny.(string)
		if !ok {
			// session value missing or corrupted
		}

		userType, err := strconv.Atoi(userTypeStr)
		if err != nil {
			// invalid role value
		}

		// Admin can access everything
		if userType == 4 {
			c.Next()
			return
		}

		// Hierarchical check
		if userType < minRole {
			utils.Render(c, 403, views.Layout(nil, views.PageData{
				Title: "Unauthorized",
			}, views.ErrorPage(messages.Error403)))
			return
		}

		c.Next()
	}
}
