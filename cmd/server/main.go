package main

import (
	// Import generated docs
	"go-rest/internal/database"
	"go-rest/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title           Inventory API
// @version         1.0
// @description     This is a sample Inventory management server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	println("Starting server...")
	// Connect to database
	database.ConnectDatabase()

	// Initialize Router
	r := gin.Default()

	// Setup Routes
	routes.SetupRoutes(r)

	// Start Server
	r.Run(":8081")
}
