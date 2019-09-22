package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const endpointCurrencies = `https://free.currconv.com/api/v7/currencies`
const endpointConvert = `https://free.currconv.com/api/v7/convert`

func init() {
	log.Println("Loading currencies..")

	currencyURL, _ := url.Parse(endpointCurrencies)
	q := currencyURL.Query()
	q.Set("apiKey", secret.CurrencyAPIKey)
	currencyURL.RawQuery = q.Encode()

	log.Printf("GET %s", currencyURL)
	resp, err := http.Get(currencyURL.String())
	if err != nil {
		log.Println("Currencies not available:", err)
		return
	}

	var response struct {
		Results map[string]struct {
			ID           string
			CurrencyName string
		}
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		log.Println("error decoding currencies response:", err)
		return
	}

	log.Println("Supported currencies:")
	for _, curr := range response.Results {
		key := strings.ToLower(curr.ID)
		// Only add currency if it does not already exist as a unit
		// doing this because we are not manually controlling the currencies available
		// like we are with the other units
		if _, ok := unitMap[key]; !ok {
			log.Println(curr.ID, curr.CurrencyName)
			curr := currencyUnit{curr.ID}
			unitMap[key] = &curr
		}
	}

	log.Println("Currencies loaded")

}

type currencyUnit struct {
	id string
}

func (cu *currencyUnit) name() string {
	return cu.id
}

func (cu *currencyUnit) fromFloat(f float64) unitVal {
	return currencyVal{f, cu}
}

type currencyVal struct {
	v float64
	u *currencyUnit
}

func (cv currencyVal) String() string {
	return fmt.Sprintf("%.2f %s", cv.v, cv.u.name())
}

func (cv currencyVal) convert(to unitType) (unitVal, error) {
	if to, ok := to.(*currencyUnit); ok {
		curr, err := convertCurrency(cv.u, to, cv.v)
		if err != nil {
			return nil, fmt.Errorf("Currency conversion not available right now")
		}
		return currencyVal{curr, to}, nil
	}
	return nil, convErr(cv.u, to)
}

func convertCurrency(from, to *currencyUnit, val float64) (float64, error) {
	convertURL, _ := url.Parse(endpointConvert)
	q := convertURL.Query()
	q.Set("apiKey", secret.CurrencyAPIKey)
	q.Set("compact", "ultra")

	op := fmt.Sprintf("%s_%s", from.id, to.id)
	q.Add("q", op)
	convertURL.RawQuery = q.Encode()

	resp, err := http.Get(convertURL.String())
	if err != nil {
		return 0, err
	}

	var response map[string]float64
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return 0, err
	}

	conversion, ok := response[op]
	if !ok {
		return 0, fmt.Errorf("Unexpected response: %v", response)
	}

	return val * conversion, nil
}
