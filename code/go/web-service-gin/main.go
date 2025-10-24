package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album 代表一张唱片专辑的相关数据
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice 用于存储唱片专辑数据
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums 响应所有唱片专辑数据的 JSON 列表
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums 添加一张唱片专辑
func postAlbums(c *gin.Context) {
	var newAlbum album

	// 绑定接收的 JSON 到 newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// 添加新的专辑到切片
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID 根据 ID 获取一张唱片专辑
func getAlbumByID(c *gin.Context) {
	// c.Param 获取 URL 路径中的参数
	id := c.Param("id")

	for _, album := range albums {
		if album.ID == id {
			// c.IndentedJSON 返回 JSON 响应
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	// 如果找不到，返回 404 错误
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)

	router.Run("localhost:8080")
}
