package session

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetSession(c *gin.Context, key string, value string) {
	sass := sessions.Default(c)
	sass.Set(key, value)
	sass.Options(sessions.Options{
		Path: "/",
	})
	if err := sass.Save(); err != nil {
		log.Println("Error saving session:", err)
	}
	log.Println("Session User data", sass.Get("user_id"))
}

func SessionInit(router *gin.Engine) {
	store := cookie.NewStore([]byte("very-secret-key"))
	router.Use(sessions.Sessions("mysession", store))
}

// AuthRequired ensures “user_id” exists in the session
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		if sess.Get("user_id") == nil {
			// not logged in → redirect or abort
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func LogoutHandler(c *gin.Context) {
	sess := sessions.Default(c)

	// 1) Remove all keys from the session
	sess.Clear()

	// 2) Tell the browser to delete the cookie immediately
	sess.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1, // <— expires the cookie right away
		// Secure, HttpOnly, SameSite, etc. can go here too
	})

	// 3) Save the session, which writes the Set‑Cookie header
	if err := sess.Save(); err != nil {
		c.String(http.StatusInternalServerError, "Failed to clear session")
		return
	}

	// 4) Redirect or render a logout confirmation
	c.Redirect(http.StatusSeeOther, "/login")
}

func GetSession(c *gin.Context, key string) string {
	sess := sessions.Default(c)
	value := sess.Get(key)

	value, ok := value.(string)
	if !ok {
		log.Printf("Session value for key '%s' is not a string: %v", key, value)
		c.Redirect(http.StatusBadRequest, "/unexpected-error")
	}
	return value.(string)
}
