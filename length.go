package main

import (
	"fmt"
	"math"

	"github.com/martinlindhe/unit"
)

type lengthUnit struct {
	unitCommon
	length unit.Length
	to     func(unit.Length) float64
}

func (lu *lengthUnit) fromFloat(f float64) unitVal {
	return lengthVal{unit.Length(f) * lu.length, lu}
}

var (
	meter      = &lengthUnit{"m", unit.Meter, unit.Length.Meters}
	kilometer  = &lengthUnit{"km", unit.Kilometer, unit.Length.Kilometers}
	millimeter = &lengthUnit{"mm", unit.Millimeter, unit.Length.Millimeters}
	centimeter = &lengthUnit{"cm", unit.Centimeter, unit.Length.Centimeters}
	nanometer  = &lengthUnit{"nm", unit.Nanometer, unit.Length.Nanometers}
	inch       = &lengthUnit{"in", unit.Inch, unit.Length.Inches}
	foot       = &lengthUnit{"ft", unit.Foot, unit.Length.Feet}
	yard       = &lengthUnit{"yd", unit.Yard, unit.Length.Yards}
	mile       = &lengthUnit{"miles", unit.Mile, unit.Length.Miles}
	furlong    = &lengthUnit{"furlongs", unit.Furlong, unit.Length.Furlongs}
	lightyear  = &lengthUnit{"ly", unit.LightYear, unit.Length.LightYears}
)

type lengthVal struct {
	v unit.Length
	u *lengthUnit
}

func (lv lengthVal) String() string {
	return simpleUnitString(lv.u.to(lv.v), lv.u)
}

func (lv lengthVal) convert(to unitType) (unitVal, error) {
	switch to := to.(type) {
	case *lengthUnit:
		lv.u = to
		return lv, nil
	case footInchUnit:
		feet, fraction := math.Modf(lv.v.Feet())
		inches := (unit.Length(fraction) * unit.Foot).Inches()
		return footInchVal{feet, inches}, nil
	default:
		return nil, convErr(lv.u, to)

	}
}

type footInchUnit struct{}

func (footInchUnit) name() string {
	return "feet+inches"
}

var footInch = footInchUnit{}

type footInchVal struct {
	feet, inches float64
}

func (val footInchVal) String() string {
	if val.inches == 0 {
		return simpleUnitString(val.feet, foot)
	}
	return fmt.Sprintf(`%.f' %.f"`, val.feet, val.inches)
}

func (val footInchVal) convert(to unitType) (unitVal, error) {
	switch to := to.(type) {
	case footInchUnit:
		return val, nil
	case *lengthUnit:
		feet := unit.Length(val.feet) * unit.Foot
		inches := unit.Length(val.inches) * unit.Inch
		return lengthVal{feet + inches, to}, nil
	default:
		return nil, convErr(footInch, to)
	}
}
