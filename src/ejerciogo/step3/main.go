package main

import (
	"ejerciogo/step2/model"

	"github.com/gin-gonic/gin"

	"log"

	"encoding/json"
	"net/http"
	"strings"

	"sync"
	"time"

	"io/ioutil"
)

func main() {
	router := gin.Default()

	router.GET("/myapi", func(c *gin.Context) {

		vCryptos := c.DefaultQuery("cryptos", "ForcePartial")

		vListCrypto := strings.Split(vCryptos, ",")

		// var vResponse model.ResponseListCotacao
		var vResponse [3]model.CotacaoMoedaResponse

		var wg sync.WaitGroup    //cria WaitGroup
		wg.Add(len(vListCrypto)) //configura qtde de goroutines
		for i, vCrypto := range vListCrypto {
			c := make(chan model.CotacaoMoedaResponse)
			go func() {
				defer wg.Done()
				time.Sleep(time.Duration(i*2) * time.Millisecond)

				callApiCrypto(vCrypto, c)
				log.Println(i)
			}()
			vResponse[i] = <-c
		}

		go func() {
			wg.Wait()
		}()

		c.JSON(http.StatusOK, vResponse)

	})

	router.Run(":8080")
}

func callApiCrypto(vCrypto string, c chan model.CotacaoMoedaResponse) {
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

	c <- vCurrency

}

//Fetch is exported ...
// func FetchCrypto(fiat string, crypto string) (string, error) {
func FetchCrypto(crypto string) (model.Cryptoresponse, error) {

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

	// Passo(01) lendo o json do response do http request e transforma em Array de Bytes
	responseJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ooopsss! erro ao criar array de bytes do Response")
		return nil, err
	}

	// Passo(02) converte array de bytes em uma struct
	var dadosJason model.Cryptoresponse
	err = json.Unmarshal(responseJson, &dadosJason)
	if err != nil {
		log.Printf("ooopsss! erro ao realizar unmarshall do response")
		return nil, err
	}

	//Invoke the text output function & return it with nil as the error value
	return dadosJason, nil
}
