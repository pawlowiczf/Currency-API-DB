package api

import (
	"currency_api/utils"
	"log"
)

// Example usage: https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/btc.json

const (
	RAW_URL = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/"
)

type DailyExchangeDetails struct {
	Date     string `json:"date"`
	Currency map[string]float64 `json:"btc"`
}
// `json:"btc"`

func GetDailyExchangeRate(fromCurrency string) (DailyExchangeDetails, error) {
	//
	var target DailyExchangeDetails = DailyExchangeDetails{}

	err := utils.MakeHTTPGetRequest(RAW_URL+fromCurrency+".json", nil, &target)
	if err != nil {
		log.Fatalln(err)
	}

	return target, err
}
