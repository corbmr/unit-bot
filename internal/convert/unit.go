package convert

import (
	"fmt"
	"strings"
)

// RegisterAliases registers aliases for a UnitType
// An alias is only applied if it does not already exist
func RegisterAliases(unit UnitType, aliases []string) {
	for _, alias := range aliases {
		if _, ok := unitMap[alias]; !ok {
			unitMap[alias] = unit
		}
	}
}

// ParseUnit parses a UnitType
func ParseUnit(s string) (UnitType, bool) {
	u, ok := unitMap[strings.ToLower(s)]
	return u, ok
}

// UnitType represent a single type of unit
type UnitType interface {
	Name() string
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
	return fmt.Sprintf("Can't convert from %s to %s", err.From.Name(), err.To.Name())
}

type unitCommon string

func (c unitCommon) Name() string {
	return string(c)
}

func simpleUnitString(f float64, u UnitType) string {
	return fmt.Sprintf("%.6g %s", f, u.Name())
}
