package model

type CotacaoMoedaResponse struct {
	Id      string              `json:"id"`
	Content ContentCotacaoMoeda `json:"content,omitempty"`
	Partial string              `json:"partial"`
}

type CotacaoMoedaResponsePartial struct {
	Id      string `json:"id"`
	Partial string `json:"partial"`
}

type ContentCotacaoMoeda struct {
	Price    string `json:"price,omitempty"`
	Currency string `json:"currency,omitempty"`
}
