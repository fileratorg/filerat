package main

import (
	"github.com/fileratorg/filerat"
	"github.com/gin-gonic/gin"
	"fmt"
	"os"
)

func main() {

	// Force debug mode
	gin.SetMode("debug")

	// Set HTTP port as environment variable
	os.Setenv("PORT", fmt.Sprintf("%d", filerat.ServerPort))

	router := gin.Default()
	router.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		age := c.DefaultQuery("age", "")

		var message string
		if age != "" {
			message = fmt.Sprintf("Hello, %s. You are %s years old.", name, age)
		} else {
			message = fmt.Sprintf("Hello, %s.", name)
		}

		c.JSON(200, gin.H{
			"message": message,
		})
	})
	router.Run()
}
