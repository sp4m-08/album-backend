package controllers

import (
	"net/http"
	"strconv"

	"rest/database"
	"rest/models"

	"github.com/gin-gonic/gin"
)

func GetAlbums(c *gin.Context) {
	var albums []models.Album
	result := database.DB.Order("id DESC").Find(&albums)

	//if result is false i.e there are no albums in database
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return //return alley
	}

	//if albums are there in db, print them
	c.IndentedJSON(http.StatusOK, albums)

}

func GetAlbumByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	var album models.Album
	result := database.DB.First(&album, id) //search for book with that id
	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}

	//if book is present then print is response
	c.IndentedJSON(http.StatusOK, album)

}

func PostAlbum(c *gin.Context) {
	var newAlbum models.Album

	//fetch request body
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	// for _, a := range models.Albums {
	// 	if a.Title == newAlbum.Title && a.Artist == newAlbum.Artist {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album already exists"})
	// 		return
	// 	}
	// }

	var existing models.Album //check if album with same details are already present
	result := database.DB.Where("title = ? AND artist= ?", newAlbum.Title, newAlbum.Artist).First(&existing)
	if result.RowsAffected > 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album already exists"})
	}

	//create album with new details
	if err := database.DB.Create(&newAlbum).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum) //send in response
}

func UpdateAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	var album models.Album
	if err := database.DB.First(&album, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	var updatedAlbum models.Album
	if err := c.BindJSON(&updatedAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	album.Title = updatedAlbum.Title
	album.Artist = updatedAlbum.Artist
	album.Price = updatedAlbum.Price

	database.DB.Save(&album)
	c.IndentedJSON(http.StatusOK, album)

	// for i, a := range models.Albums {
	// 	if a.ID == id {
	// 		updatedAlbum.ID = id
	// 		models.Albums[i] = updatedAlbum
	// 		c.IndentedJSON(http.StatusOK, updatedAlbum)
	// 		return
	// 	}
	// }
	// c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func DeleteAlbumByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	var album models.Album
	if err := database.DB.First(&album, id).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	//if album is fetched, delete it from db
	database.DB.Delete(&album)
	c.IndentedJSON(http.StatusOK, gin.H{"message ": "Album deleted"})

	// for i, a := range albums {
	// 	if a.ID == id {
	// 		models.Albums = append(models.Albums[:i], models.Albums[i+1:]...)
	// 		for idx := range models.Albums {
	// 			models.Albums[idx].ID = idx + 1
	// 		}
	// 		c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted album successfully"})
	// 		return
	// 	}
	// }
	// c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
