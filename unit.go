package convert

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

var (
	unitAliasMap     map[string]UnitType
	unitDimensionMap map[UnitDimension][]UnitType
	unitLock         sync.RWMutex
)

func init() {
	unitAliasMap = make(map[string]UnitType)
	unitDimensionMap = make(map[UnitDimension][]UnitType)
	refreshUnitMaps()
}

func refreshUnitMaps() {
	for unit, aliases := range supportedUnits {
		for _, alias := range aliases {
			alias = strings.ToLower(alias)
			if _, ok := unitAliasMap[alias]; !ok {
				unitAliasMap[alias] = unit
			}
		}
		dim := unit.Dimension()
		unitDimensionMap[dim] = append(unitDimensionMap[dim], unit)
	}

	slog.Info("Refreshed unit maps",
		"unitAliasMap", unitAliasMap, "unitDimensionMap", unitDimensionMap)
}

// LookupUnit parses a UnitType.
// Lazily loads currency units
func LookupUnit(s string) (UnitType, bool) {
	s = strings.ToLower(s)
	unitLock.RLock()
	defer unitLock.RUnlock()
	u, ok := unitAliasMap[s]
	if !ok {
		unitLock.RUnlock()
		currencyOnce.Do(loadCurrencies)
		unitLock.RLock()
		u, ok = unitAliasMap[s]
	}
	return u, ok
}

type UnitDimension int

const (
	UnitDimensionNone UnitDimension = iota
	UnitDimensionLength
	UnitDimensionMass
	UnitDimensionSpeed
	UnitDimensionDuration
	UnitDimensionTemperature
	UnitDimensionVolume
	UnitDimensionCurrency
)

// UnitType represent a single type of unit
type UnitType interface {
	fmt.Stringer
	Dimension() UnitDimension
	FromFloat(float64) UnitVal
}

// UnitVal is a value with unit that can be converted to another unit
type UnitVal interface {
	fmt.Stringer
	Unit() UnitType
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
	dimension UnitDimension
	name      string
	fromFloat func(float64) U
	toFloat   func(U) float64
}

func (u *SimpleUnit[U]) FromFloat(f float64) UnitVal {
	return SimpleUnitValue[U]{
		value: u.fromFloat(f),
		unit:  u,
	}
}

func (u *SimpleUnit[U]) Dimension() UnitDimension {
	return u.dimension
}

func (u *SimpleUnit[U]) String() string {
	return u.name
}

type SimpleUnitValue[U ~float64] struct {
	value U
	unit  *SimpleUnit[U]
}

func (v SimpleUnitValue[U]) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*SimpleUnit[U]); ok {
		v.unit = to
		return v, nil
	}
	return nil, ErrorConversion{v.unit, to}
}

func (v SimpleUnitValue[U]) Unit() UnitType {
	return v.unit
}

func (v SimpleUnitValue[U]) String() string {
	return simpleUnitString(v.unit.toFloat(v.value), v.unit)
}
