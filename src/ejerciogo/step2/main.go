package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()

	router.GET("/myapi", func(c *gin.Context) {
		data := c.DefaultQuery("data","")

		resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
		if err != nil {
		   log.Fatalln(err)
		}
	 //We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		   log.Fatalln(err)
		}
	 //Convert the body to type string
		sb := string(body)
		
		// c.STRING
		
		// type StructRespJson struct {
		// 	userid string `json:"userId"`
		// 	id string `json:"id"`
		// 	title string `json:"title"`
		// 	body string `json:"body"`
		// }

		// var respbody StructRespJson{
		
		// }

		// responsejson := json.Marshal()

		c.STRING(200, sb)

		// c.JSON(200, responsejson)


		// var dataresponsejson = DataResponseJson{data}

		// c.JSON(200, dataresponsejson)

	})

	router.Run(":8080")

	}

	func HttpGet (url string) *http.Response{
		resp, err := http.Get(url)
		if err != nil {
		   log.Fatalln(err)
		}

	 	return resp
	}