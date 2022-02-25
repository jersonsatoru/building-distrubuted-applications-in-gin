package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", helloWorldHandler)
	r.Run()
}

func helloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello World",
	})
}
