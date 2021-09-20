package main

import (
	"github.com/gin-gonic/gin"
	model "go.mod/internal"
	"log"
	"sync"

	"encoding/json"
	"net/http"
	"strings"

	"io/ioutil"
)





func main() {

	router := gin.Default()
	router.GET("/myapi", func(c *gin.Context){
		processaget(c)
	})

	router.Run(":8080")
}

func processaget(c *gin.Context){

	vListCrypto, flagconcurrency := getparameters(c)

	var vResponse []model.CotacaoMoedaResponse

	if flagconcurrency == "S" {
		// Executa com Goroutines
		vResp := processwithconcurrency(vListCrypto)
		for _, v := range vResp {
			vResponse = append(vResponse, v)
		}
	} else {
		// Executa sem Goroutines
		vResp := processwithoutconcurrency(vListCrypto)
		for _, v := range vResp {
			vResponse = append(vResponse, v)
		}
	}

	returnresponse(vResponse, c)
}

func returnresponse(vResponse []model.CotacaoMoedaResponse, c *gin.Context)  {
	if haspartial(vResponse){
		c.JSON(http.StatusPartialContent, vResponse)
	} else {
		c.JSON(http.StatusOK, vResponse)
	}
}

func getparameters(c *gin.Context) ([]string, string){
	vCryptos := c.DefaultQuery("cryptos", "ForcePartial")
	flagconcurrency := strings.ToUpper(c.DefaultQuery("concurrency", "S"))

	vListCrypto := strings.Split(vCryptos, ",")

	return vListCrypto, flagconcurrency
}

func processwithconcurrency(vListCrypto []string,) ([]model.CotacaoMoedaResponse){
	var vResponse []model.CotacaoMoedaResponse
	var wg sync.WaitGroup    //cria WaitGroup
	wg.Add(len(vListCrypto)) //configura qtde de goroutines
	cCurrency := make(chan model.CotacaoMoedaResponse, len(vListCrypto))

	//for _, vCrypto := range vListCrypto {
	for _, vCrypto := range vListCrypto {
		log.Println(vCrypto)
		go func() {
			defer wg.Done()
			callApiCrypto(vCrypto, cCurrency)
		}()
	}

	wg.Wait()

	for i := 0; i < len(vListCrypto); i++ {
		vResp := <-cCurrency
		vResponse = append(vResponse, vResp)
	}

	return vResponse
}

func processwithoutconcurrency(vListCrypto []string) ([]model.CotacaoMoedaResponse){
	var vResponse []model.CotacaoMoedaResponse
	for _, vCrypto := range vListCrypto {
		var vResp model.CotacaoMoedaResponse
		vResp, _ = callApiCryptoNoConcurrency(vCrypto)
		vResponse = append(vResponse, vResp)
	}
	return vResponse
}

func haspartial(vResponse []model.CotacaoMoedaResponse) (bool) {
	for _, v := range vResponse {
		if v.HasPartial() {
			return true
		}
	}
	return false
}


func callApiCrypto(vCrypto string, cCurrency chan model.CotacaoMoedaResponse) {

	crypto, err := FetchCrypto(vCrypto)

	var vCurrency model.CotacaoMoedaResponse
	if err != nil {
		log.Println(err)
		vCurrency = model.CotacaoMoedaResponse{Id: vCrypto, Partial: "true"}

	} else {
		list := make(model.CryptoResponse, 0, len(crypto))
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

		list := make(model.CryptoResponse, 0, len(crypto))
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

func FetchCrypto(crypto string) (model.CryptoResponse, error) {

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
	var dadosJason model.CryptoResponse
	err = json.Unmarshal(responseJson, &dadosJason)
	if err != nil {
		log.Printf("ooopsss! erro ao realizar unmarshall do response")
		return nil, err
	}

	return dadosJason, nil
}

