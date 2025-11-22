package main

import (
	"go-rest/internal/database"
	"go-rest/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.ConnectDatabase()

	// Initialize Router
	r := gin.Default()

	// Setup Routes
	routes.SetupRoutes(r)

	// Start Server
	r.Run(":8081")
}
