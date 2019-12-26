package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-snippets/go-mongo-crud/app"
)

func InitServer() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/", app.AddUser)
		v1.GET("/", app.GetAllUsers)
		v1.GET("/:id", app.GetUser)
		v1.PUT("/:id", app.UpdateUser)
		v1.DELETE("/:id", app.RemoveUser)
	}

	v2 := router.Group("/api/v2")
	{
		v2.POST("/date", timeFormatter)
	}

	v3 := router.Group("/api/v3")
	{
		v3.POST("/perlin/2D", app.GetPerlinNoise2D)
	}

	router.Run()
}

func timeFormatter(c *gin.Context) {
	datetime := c.PostForm("date")

	layout := "2006-01-02T15:04:05"

	t, _ := time.Parse(layout, datetime)
	fmt.Println(t.Format("2006-01-02T15:04:05"))
	if t.Hour() < 12 {
		t = t.Add(time.Hour * time.Duration(12))
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"date": t.Format("2006-01-02T15:04:05")})
}
