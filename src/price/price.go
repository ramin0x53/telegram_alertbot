package price

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type priceResBody struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}

// Get price of binance
func GetPrice(symbol string) string {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=" + strings.ToUpper(symbol))
	if err != nil {
		log.Fatalln(err)
	}

	body := &priceResBody{}
	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return ""
	}

	if body.Price == "" {
		return ""
	}

	return body.Price
}
