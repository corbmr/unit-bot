package main

import "fmt"

type unitType interface {
	name() string
}

type unitCommon string

func (c unitCommon) name() string {
	return string(c)
}

type unitVal interface {
	convert(to unitType) (unitVal, error)
	fmt.Stringer
}

type simpleUnit interface {
	fromFloat(f float64) unitVal
}

func convErr(from, to unitType) error {
	return fmt.Errorf("Can't convert from %s to %s", from.name(), to.name())
}

func simpleUnitString(f float64, u unitType) string {
	return fmt.Sprintf("%.6g %s", f, u.name())
}
