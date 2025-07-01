package main

import (
	"log"
	"rest/database"
	"rest/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connect()
	router := gin.Default()
	routes.RegisterAlbumRoutes(router)
	router.Run("localhost:9090")
}
