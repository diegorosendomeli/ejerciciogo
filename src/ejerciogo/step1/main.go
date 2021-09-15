package main

import (
	"github.com/gin-gonic/gin"

	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/myapi", func(c *gin.Context) {
		data := c.DefaultQuery("data","")

		type DataResponseJson struct {
			Data string `json:"data"`
		}

		var dataresponsejson = DataResponseJson{data}

		c.JSON(200, dataresponsejson)

	})

	router.Run(":8080")

	}