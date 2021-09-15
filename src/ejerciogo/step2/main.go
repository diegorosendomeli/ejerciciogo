package main

import (
	"ejerciogo/step2/model"

	"github.com/gin-gonic/gin"

	"log"

	"encoding/json"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/myapi", func(c *gin.Context) {

		vCrypto := c.DefaultQuery("crypto", "ForcePartial")

		crypto, err := FetchCrypto(vCrypto)

		if err != nil {
			log.Println(err)
			vCurrencyPartial := model.CotacaoMoedaResponsePartial{Id: vCrypto, Partial: "true"}
			c.JSON(206, vCurrencyPartial)
		} else {

			list := make(model.Cryptoresponse, 0, len(crypto))
			list = append(list, crypto...)

			var vCurrency model.CotacaoMoedaResponse
			for _, v := range list {
				vCurrency = model.CotacaoMoedaResponse{
					Id: v.Id,
					Content: model.ContentCotacaoMoeda{
						Price:    v.Price,
						Currency: v.Currency,
					},
					Partial: "false",
				}
			}

			if vCurrency.Id == "" {
				log.Println("Content Empty - Return Partial")
				vCurrencyPartial := model.CotacaoMoedaResponsePartial{Id: vCrypto, Partial: "true"}
				c.JSON(206, vCurrencyPartial)
			} else {
				c.JSON(200, vCurrency)
			}

		}

	})

	router.Run(":8080")
}

//Fetch is exported ...
// func FetchCrypto(fiat string, crypto string) (string, error) {
func FetchCrypto(crypto string) (model.Cryptoresponse, error) {

	// tempCrypto := "BTC"

	//Build The URL string
	URL := "https://api.nomics.com/v1/currencies/ticker?key=3990ec554a414b59dd85d29b2286dd85&interval=1d&ids=" + crypto
	// URL := "https://api.nomics.com/v1/currencies/ticker?key=3990ec554a414b59dd85d29b2286dd85&interval=1d&ids=" + tempCrypto

	log.Printf(URL)

	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("ooopsss! erro na chamada da API")
		return nil, err
	}
	defer resp.Body.Close()
	//Create a variable of the same type as our model
	var cResp model.Cryptoresponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		log.Printf("ooopsss! erro no decode")
		return nil, err
	}
	//Invoke the text output function & return it with nil as the error value
	return cResp, nil
}
