package convert

import (
	"fmt"
	"strings"
)

// UnitType represent a single type of unit
type UnitType interface {
	Name() string
}

type unitCommon string

func (c unitCommon) Name() string {
	return string(c)
}

// UnitVal is a value with unit
type UnitVal interface {
	fmt.Stringer
	Convert(to UnitType) (UnitVal, error)
}

// SimpleUnit is a UnitType that can be created from float conversion
// This is most units
type SimpleUnit interface {
	UnitType

	FromFloat(float64) UnitVal
}

func convErr(from, to UnitType) error {
	return fmt.Errorf("Can't convert from %s to %s", from.Name(), to.Name())
}

func simpleUnitString(f float64, u SimpleUnit) string {
	return fmt.Sprintf("%.6g %s", f, u.Name())
}

// ParseUnit parses a UnitType
func ParseUnit(s string) (UnitType, error) {
	u, ok := unitMap[strings.ToLower(s)]
	if !ok {
		return nil, fmt.Errorf("Invalid unit %s", s)
	}
	return u, nil
}
