package main

import (
	"rest/database"
	"rest/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	router := gin.Default()
	routes.RegisterAlbumRoutes(router)
	router.Run("localhost:9090")
}
