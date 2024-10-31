package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!\n")
}

func main() {
	r := gin.Default()
	r.GET("/hello", helloWorld)
	r.Run()
}
