package main

import (
	"gl/db"
	"gl/routes"
	"gl/session"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	session.SessionInit(r) // Initialize session management
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; will rely on actual environment variables")
	}
	db.Init()
	routes.RegisterRoutes(r)
	r.Run(":8080") // Start the server on port 8080

}
