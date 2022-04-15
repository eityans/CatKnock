package main

import (
	"github.com/gin-gonic/gin"
)

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
			"message": "ここをキャンプ地とする",
	})
}

func callbackHandler(c *gin.Context) {
	c.JSON(200, gin.H{
			"message": "aaa",
	})
}

func main() {
	router := gin.Default()
	router.GET("/ping", pingHandler)
	router.GET("/callback", callbackHandler)
	router.Run()
}
