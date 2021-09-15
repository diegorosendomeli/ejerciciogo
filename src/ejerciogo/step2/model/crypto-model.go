package model

// Cryptoresponse is exported, it models the data we receive.
type Cryptoresponse []struct {
	Id       string `json:"id"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
}
