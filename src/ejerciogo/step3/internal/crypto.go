package model

// Cryptoresponse is exported, it models the data we receive.
type CryptoResponse []struct {
	Id       string `json:"id"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
}

type CotacaoMoedaResponse struct {
	Id      string               `json:"id"`
	Content *ContentCotacaoMoeda `json:"content,omitempty"`
	Partial string               `json:"partial"`
}

type ContentCotacaoMoeda struct {
	Price    string `json:"price,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type ResponseListCotacao struct {
	Cotacoes [3]CotacaoMoedaResponse `json:",omitname"`
}

func (s CotacaoMoedaResponse) HasPartial() bool {
		if s.Partial == "true" {
			return true
		}
	return false
}