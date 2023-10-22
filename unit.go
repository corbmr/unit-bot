package convert

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// InitCurrency registers a callback to retrieve the apiKey for registering currencies
	// Currencies are lazily loaded only when needed
	CurrencyInit func() (string, error)
)

var (
	unitMap  map[string]UnitType
	unitLock sync.RWMutex

	currencyApiKey string
	currencyOnce   sync.Once
	currencyCache  *cache.Cache
)

func init() {
	unitMap = make(map[string]UnitType)
	refreshUnitMap()
	currencyCache = cache.New(24*time.Hour, 1*time.Hour)
}

func refreshUnitMap() {
	for unit, aliases := range supportedUnits {
		for _, alias := range aliases {
			alias = strings.ToLower(alias)
			if _, ok := unitMap[alias]; !ok {
				unitMap[alias] = unit
			}
		}
	}
}

// SupportedUnits returns all of the supported unit types mapped to a list of alises
func SupportedUnits() map[UnitType][]string {
	unitLock.RLock()
	defer unitLock.RUnlock()
	supported := make(map[UnitType][]string, len(supportedUnits))
	for unit, aliases := range supportedUnits {
		supported[unit] = aliases
	}
	return supported
}

// ParseUnit parses a UnitType.
// Lazily loads currency units
func ParseUnit(s string) (UnitType, bool) {
	s = strings.ToLower(s)
	unitLock.RLock()
	defer unitLock.RUnlock()
	u, ok := unitMap[s]
	if !ok {
		unitLock.RUnlock()
		currencyOnce.Do(loadCurrencies)
		unitLock.RLock()
		u, ok = unitMap[s]
	}
	return u, ok
}

// UnitType represent a single type of unit
type UnitType interface {
	fmt.Stringer
	FromFloat(float64) UnitVal
}

// UnitVal is a value with unit that can be converted to another unit
type UnitVal interface {
	fmt.Stringer
	Convert(to UnitType) (UnitVal, error)
}

// ErrorConversion occurs when a UnitType cannot be converted to another UnitType
type ErrorConversion struct {
	From, To UnitType
}

func (err ErrorConversion) Error() string {
	return fmt.Sprintf("Can't convert from %s to %s", err.From.String(), err.To.String())
}

func simpleUnitString(f float64, u UnitType) string {
	return fmt.Sprintf("%.6g %s", f, u.String())
}

func from[U ~float64](base U) func(float64) U {
	return func(f float64) U {
		return U(f) * base
	}
}

type SimpleUnit[U ~float64] struct {
	name      string
	fromFloat func(float64) U
	toFloat   func(U) float64
}

func (u *SimpleUnit[U]) FromFloat(f float64) UnitVal {
	return SimpleUnitValue[U]{
		V: u.fromFloat(f),
		U: u,
	}
}

func (u *SimpleUnit[U]) String() string {
	return u.name
}

type SimpleUnitValue[U ~float64] struct {
	V U
	U *SimpleUnit[U]
}

func (v SimpleUnitValue[U]) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*SimpleUnit[U]); ok {
		v.U = to
		return v, nil
	}
	return nil, ErrorConversion{v.U, to}
}

func (v SimpleUnitValue[U]) String() string {
	return simpleUnitString(v.U.toFloat(v.V), v.U)
}
