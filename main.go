package main

import (
	"fmt"
	"net/http"

	"github.com/Darrell-Devana/be-note-app/util"
	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!")
}

func main() {
	message := "be-note-app started"
	messageUpper := util.ToUpperCase(message)
	fmt.Println(messageUpper)

	r := gin.Default()
	r.GET("/hello", helloWorld)
	r.Run()
}
