package main

import (
	"ejerciogo/step2/model"
	"time"

	"github.com/gin-gonic/gin"

	"log"

	"encoding/json"
	"net/http"
	"strings"

	"sync"

	"io/ioutil"
)

func main() {
	router := gin.Default()

	router.GET("/myapi", func(c *gin.Context) {

		vCryptos := c.DefaultQuery("cryptos", "ForcePartial")
		flagconcurrency := strings.ToUpper(c.DefaultQuery("concurrency", "S"))

		vListCrypto := strings.Split(vCryptos, ",")
		var flagPartial bool = false
		var vResponse [3]model.CotacaoMoedaResponse

		if flagconcurrency == "S" {

			// Executa com Goroutines

			var wg sync.WaitGroup    //cria WaitGroup
			wg.Add(len(vListCrypto)) //configura qtde de goroutines
			for i, vCrypto := range vListCrypto {

				cCurrency := make(chan model.CotacaoMoedaResponse)
				go func() {
					defer wg.Done()

					callApiCrypto(vCrypto, cCurrency)

					log.Println(i)
				}()
				vResponse[i] = <-cCurrency
				if vResponse[i].Partial == "true" {
					flagPartial = true
				}
			}

			go func() {
				wg.Wait()
			}()

		} else {
			// Executa sem Goroutines

			for i, vCrypto := range vListCrypto {
				time.Sleep(time.Duration(i*2) * time.Millisecond)
				vResponse[i], _ = callApiCryptoNoConcurrency(vCrypto)
				if vResponse[i].Partial == "true" {
					flagPartial = true
				}
			}

		}

		if flagPartial {
			c.JSON(http.StatusPartialContent, vResponse)
		} else {
			c.JSON(http.StatusOK, vResponse)
		}

	})

	router.Run(":8080")
}

func callApiCrypto(vCrypto string, cCurrency chan model.CotacaoMoedaResponse) {

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
				vCurrency = model.CotacaoMoedaResponse{
					Id: v.Id,
					Content: &model.ContentCotacaoMoeda{
						Price:    v.Price,
						Currency: v.Currency,
					},
					Partial: "false",
				}

			} else {
				vCurrency = model.CotacaoMoedaResponse{Id: vCrypto, Partial: "true"}
			}
		}

	}
	cCurrency <- vCurrency
}

func callApiCryptoNoConcurrency(vCrypto string) (model.CotacaoMoedaResponse, error) {
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
				vCurrency = model.CotacaoMoedaResponse{
					Id: v.Id,
					Content: &model.ContentCotacaoMoeda{
						Price:    v.Price,
						Currency: v.Currency,
					},
					Partial: "false",
				}

			} else {
				vCurrency = model.CotacaoMoedaResponse{Id: vCrypto, Partial: "true"}
			}
		}
	}

	return vCurrency, nil

}

func FetchCrypto(crypto string) (model.Cryptoresponse, error) {

	URL := "https://api.nomics.com/v1/currencies/ticker?key=3990ec554a414b59dd85d29b2286dd85&interval=1d&ids=" + crypto

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

	return dadosJason, nil
}
