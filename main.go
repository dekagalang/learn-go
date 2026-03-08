package main

import (
	"log"

	"crud-api/database"
	"crud-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	database.Connect()

	// Run Atlas-based migrations
	if err := database.RunMigrations(database.SqlDB); err != nil {
		log.Fatal("Migration failed:", err)
	}

	api := r.Group("/api")

	{
		// Users routes
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)
		api.POST("/users", handlers.CreateUser)
		api.PUT("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)

		// Products routes
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.POST("/products", handlers.CreateProduct)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)
	}

	r.Run(":8080")
}
