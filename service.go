package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// define structs

// setup router
func setupRouter() *gin.Engine {
	r := gin.Default()
	// define routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}

func main() {
	r := setupRouter()
	// run server
	r.Run(":8080")
}
