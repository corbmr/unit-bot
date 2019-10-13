package convert

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

var apiKey string

// InitCurrency registers available currencies
func InitCurrency(key string) {
	apiKey = key
	loadCurrencies()
}

func loadCurrencies() {
	log.Println("Loading currencies..")

	currencyURL, _ := url.Parse(endpointCurrencies)
	q := currencyURL.Query()
	q.Set("apiKey", apiKey)
	currencyURL.RawQuery = q.Encode()

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
			curr := CurrencyUnit{curr.ID}
			unitMap[key] = &curr
		}
	}

	log.Println("Currencies loaded")
}

// CurrencyUnit is a unit of currency
type CurrencyUnit struct {
	id string
}

// Name implements UnitType for Currency
func (cu *CurrencyUnit) Name() string {
	return cu.id
}

// FromFloat implements SimpleUnit
func (cu *CurrencyUnit) FromFloat(f float64) UnitVal {
	return CurrencyVal{f, cu}
}

// CurrencyVal is a currency value with unit
type CurrencyVal struct {
	V float64
	U *CurrencyUnit
}

func (cv CurrencyVal) String() string {
	return fmt.Sprintf("%.2f %s", cv.V, cv.U.Name())
}

// Convert implements UnitVal conversion
func (cv CurrencyVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*CurrencyUnit); ok {
		curr, err := convertCurrency(cv.U, to, cv.V)
		if err != nil {
			return nil, fmt.Errorf("Currency conversion not available right now")
		}
		return CurrencyVal{curr, to}, nil
	}
	return nil, convErr(cv.U, to)
}

func convertCurrency(from, to *CurrencyUnit, val float64) (float64, error) {
	convertURL, _ := url.Parse(endpointConvert)
	q := convertURL.Query()
	q.Set("apiKey", apiKey)
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
