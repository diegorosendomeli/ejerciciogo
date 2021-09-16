package main

import (
	"ejerciogo/step2/model"

	"github.com/gin-gonic/gin"

	"log"

	"encoding/json"
	"net/http"
	"strings"
)

func main() {
	router := gin.Default()

	router.GET("/myapi", func(c *gin.Context) {

		vCryptos := c.DefaultQuery("cryptos", "ForcePartial")

		vListCrypto := strings.Split(vCryptos, "")

		// var vResponse model.ResponseListCotacao
		var vResponse [3]model.CotacaoMoedaResponse

		for i, vCrypto := range vListCrypto {
			vCurrency, _ := callApiCrypto(vCrypto)
			vResponse[i] = vCurrency

			//// PASSEI PRA FUNC E PAROU DE FUNCIONAR... DANDO TIMEOUT... VERIFICAR
			///// ADICIONAR CONCONRRENCIA

		}

		c.JSON(http.StatusOK, vResponse)

	})

	router.Run(":8080")
}

func callApiCrypto(vCrypto string) (model.CotacaoMoedaResponse, error) {
	crypto, err := FetchCrypto(vCrypto)

	var vCurrency model.CotacaoMoedaResponse
	if err != nil {
		log.Println(err)
		vCurrency = model.CotacaoMoedaResponse{Id: vCrypto, Partial: "true"}

	} else {

		list := make(model.Cryptoresponse, 0, len(crypto))
		list = append(list, crypto...)

		for _, v := range list {

			if v.Id != "" {
				//partial=false
				vCurrency = model.CotacaoMoedaResponse{
					Id: v.Id,
					Content: &model.ContentCotacaoMoeda{
						Price:    v.Price,
						Currency: v.Currency,
					},
					Partial: "false",
				}

			} else {
				//partial=true
				vCurrency = model.CotacaoMoedaResponse{Id: vCrypto, Partial: "true"}
			}
		}

	}

	return vCurrency, nil

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
