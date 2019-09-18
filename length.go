package main

import "github.com/martinlindhe/unit"

type lengthUnit struct {
	unitCommon
	length unit.Length
	to     func(unit.Length) float64
}

func (from *lengthUnit) convert(f float64, to *lengthUnit) float64 {
	return to.to(unit.Length(f) * from.length)
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

func convFeetInches(f, i float64, to *lengthUnit) float64 {
	from := unit.Length(f)*unit.Foot + unit.Length(i)*unit.Inch
	return to.to(from)
}
