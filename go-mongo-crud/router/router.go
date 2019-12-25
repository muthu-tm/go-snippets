package router

import (
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
	router.Run()
}
