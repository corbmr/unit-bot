package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	endpoint           = "https://free.currconv.com/api/v7"
	endpointCurrencies = endpoint + "/currencies"
	endpointConvert    = endpoint + "/convert"
)

var (
	// ErrorCurrencyService occurs when there is an error while calling the currency service
	ErrorCurrencyService = errors.New("currency conversion not available right now")

	currencyApiKey string
	currencyOnce   sync.Once
	currencyCache  *cache.Cache = cache.New(24*time.Hour, 1*time.Hour)

	extraAliases = map[string][]string{
		"USD": {"$", "US$", "dollar", "dollars"},
		"EUR": {"€", "euro", "euros"},
		"JPY": {"¥", "yen"},
		"GBP": {"£"},
		"CAD": {"CA$"},
		"AUS": {"AUS$"},
	}
)

func SetCurrencyApiKey(apiKey string) {
	currencyApiKey = apiKey
}

type supportedCurrencies struct {
	Results map[string]struct {
		ID           string
		CurrencyName string
	}
}

func loadCurrencies() {
	if currencyApiKey == "" {
		slog.Info("Currency API key was not set. Currency conversion is not available")
		return
	}
	slog.Info("Loading currencies..")

	currencies, err := retrieveSupportedCurrencies()
	if err != nil {
		slog.Error("Error loading currencies:", err)
		return
	}

	slog.Info("Supported currencies:")
	unitLock.Lock()
	defer unitLock.Unlock()
	for _, curr := range currencies.Results {
		slog.Info(curr.ID, curr.CurrencyName)
		unit := &CurrencyUnit{curr.ID}
		supportedUnits[unit] = append(supportedUnits[unit], curr.ID)
		if aliases, ok := extraAliases[unit.id]; ok {
			supportedUnits[unit] = append(supportedUnits[unit], aliases...)
		}
	}
	refreshUnitMaps()

	slog.Info("Currencies loaded")
}

func retrieveSupportedCurrencies() (*supportedCurrencies, error) {
	currencyURL, _ := url.Parse(endpointCurrencies)
	q := currencyURL.Query()
	q.Set("apiKey", currencyApiKey)
	currencyURL.RawQuery = q.Encode()

	resp, err := http.Get(currencyURL.String())
	if err != nil {
		return nil, fmt.Errorf("Error calling currency service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	var response supportedCurrencies
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode currencies response: %w\nBody: %v", err, string(body))
	}

	return &response, nil
}

// CurrencyUnit is a unit of currency
type CurrencyUnit struct {
	id string
}

// Name implements UnitType for Currency
func (cu *CurrencyUnit) String() string {
	return cu.id
}

// FromFloat implements SimpleUnit
func (cu *CurrencyUnit) FromFloat(f float64) UnitVal {
	return CurrencyVal{f, cu}
}

func (cu *CurrencyUnit) Dimension() UnitDimension {
	return UnitDimensionCurrency
}

// CurrencyVal is a currency value with unit
type CurrencyVal struct {
	V float64
	U *CurrencyUnit
}

func (cv CurrencyVal) String() string {
	return fmt.Sprintf("%.2f %s", cv.V, cv.U)
}

// Convert implements UnitVal conversion
func (cv CurrencyVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*CurrencyUnit); ok {
		rate, err := getRate(cv.U, to)
		if err != nil {
			if err != ErrorCurrencyService {
				slog.Error("Error calling currency service,", err)
			}
			return nil, ErrorCurrencyService
		}
		return CurrencyVal{cv.V * rate, to}, nil
	}
	return nil, ErrorConversion{cv.U, to}
}

func (cv CurrencyVal) Unit() UnitType {
	return cv.U
}

func getRate(from, to *CurrencyUnit) (float64, error) {
	op := from.id + "_" + to.id
	rate, ok := currencyCache.Get(op)
	if ok {
		slog.Debug("Cache hit:", op)
		return rate.(float64), nil
	} else {
		slog.Debug("Cache miss:", op)
		r, err := getRateNoCache(op)
		if err != nil {
			return 0, err
		}
		currencyCache.Set(op, r, cache.DefaultExpiration)
		return r, nil
	}
}

func getRateNoCache(op string) (float64, error) {
	if len(currencyApiKey) == 0 {
		return 0, ErrorCurrencyService
	}

	convertURL, _ := url.Parse(endpointConvert)
	q := convertURL.Query()
	q.Set("apiKey", currencyApiKey)
	q.Set("compact", "ultra")
	q.Set("q", op)
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
		return 0, fmt.Errorf("unexpected response: %v", response)
	}

	return conversion, nil
}
