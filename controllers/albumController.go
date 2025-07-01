package controllers

import (
	"net/http"
	"strconv"

	"rest/database"
	"rest/models"
	"rest/utils"

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

	title := c.PostForm("title")
	artist := c.PostForm("artist")
	priceStr := c.PostForm("price")

	price, err := strconv.ParseFloat(priceStr, 64)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid Price"})
		return
	}
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	imageURL, err := utils.UploadFileToS3(file, fileHeader, title+"-"+artist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var existingAlbum models.Album
	result := database.DB.Where("title =? AND artist=?", title, artist).First(&existingAlbum)
	if result.RowsAffected > 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "album already exists"})
		return
	}

	newAlbum := models.Album{
		Title:    title,
		Artist:   artist,
		Price:    price,
		ImageURL: imageURL,
	} //since we using form body

	if err := database.DB.Create(&newAlbum).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, newAlbum)

	//below is without image context
	// var newAlbum models.Album

	// //fetch request body
	// if err := c.BindJSON(&newAlbum); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
	// 	return
	// }

	// // for _, a := range models.Albums {
	// // 	if a.Title == newAlbum.Title && a.Artist == newAlbum.Artist {
	// // 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album already exists"})
	// // 		return
	// // 	}
	// // }

	// var existing models.Album //check if album with same details are already present
	// result := database.DB.Where("title = ? AND artist= ?", newAlbum.Title, newAlbum.Artist).First(&existing)
	// if result.RowsAffected > 0 {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album already exists"})
	// }

	// //create album with new details
	// if err := database.DB.Create(&newAlbum).Error; err != nil {
	// 	c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.IndentedJSON(http.StatusCreated, newAlbum) //send in response
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

	// var updatedAlbum models.Album //for json body, now we are sending form body
	// if err := c.BindJSON(&updatedAlbum); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
	// 	return
	// }

	title := c.PostForm("title")
	artist := c.PostForm("artist")
	priceStr := c.PostForm("price") //form data parsing

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	album.Title = title
	album.Artist = artist
	album.Price = price

	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded image"})
			return
		}
		defer file.Close()

		imageURL, err := utils.UploadFileToS3(file, fileHeader, title+"-"+artist)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image"})
			return
		}
		album.ImageURL = imageURL
	}

	if err := database.DB.Save(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
