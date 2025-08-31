package main

import (
	"fmt"
	"go-ticketing/config"
	"go-ticketing/models"
	"go-ticketing/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func InitializeApp() *gin.Engine {
	r := gin.Default()

	db := config.ConnectDatabase()

	// Auto migrate tabel
	db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Ticket{},
	)

	routes.SetupRoutes(r, db)

	return r
}

func main() {
	// Normal start server
	app := InitializeApp()
	//membaca environment variable yang sudah ada di proses.
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // fallback
	}
	fmt.Println("Server running on port", port)
	app.Run(":" + port)
}
