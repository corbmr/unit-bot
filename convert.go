package main

import (
	"fmt"
	"strings"

	"github.com/martinlindhe/unit"
)

type unitType interface {
	name() string
}

type unitCommon string

func (c unitCommon) name() string {
	return string(c)
}

type lengthUnit struct {
	unitCommon
	length unit.Length
	to     func(unit.Length) float64
}

func (from *lengthUnit) convert(f float64, to *lengthUnit) float64 {
	return to.to(unit.Length(f) * from.length)
}

type weightUnit struct {
	unitCommon
	weight unit.Mass
	to     func(unit.Mass) float64
}

func (from *weightUnit) convert(f float64, to *weightUnit) float64 {
	return to.to(unit.Mass(f) * from.weight)
}

type tempertureUnit struct {
	unitCommon
	from func(float64) unit.Temperature
	to   func(unit.Temperature) float64
}

func (from *tempertureUnit) convert(f float64, to *tempertureUnit) float64 {
	return to.to(from.from(f))
}

var (
	meter      = &lengthUnit{"m", unit.Meter, unit.Length.Meters}
	kilometer  = &lengthUnit{"km", unit.Kilometer, unit.Length.Kilometers}
	millimeter = &lengthUnit{"mm", unit.Millimeter, unit.Length.Millimeters}
	centimeter = &lengthUnit{"cm", unit.Centimeter, unit.Length.Centimeters}
	inch       = &lengthUnit{"in", unit.Inch, unit.Length.Inches}
	foot       = &lengthUnit{"ft", unit.Foot, unit.Length.Feet}
	mile       = &lengthUnit{"miles", unit.Mile, unit.Length.Miles}
	furlong    = &lengthUnit{"furlongs", unit.Furlong, unit.Length.Furlongs}

	gram     = &weightUnit{"g", unit.Gram, unit.Mass.Grams}
	kilogram = &weightUnit{"kg", unit.Kilogram, unit.Mass.Kilograms}
	pound    = &weightUnit{"lbs", unit.AvoirdupoisPound, unit.Mass.AvoirdupoisPounds}

	celcius    = &tempertureUnit{"°C", unit.FromCelsius, unit.Temperature.Celsius}
	fahrenheit = &tempertureUnit{"°F", unit.FromFahrenheit, unit.Temperature.Fahrenheit}
	kelvin     = &tempertureUnit{"K", unit.FromKelvin, unit.Temperature.Kelvin}
)

var unitMap = map[string]unitType{
	"m":           meter,
	"meter":       meter,
	"meters":      meter,
	"km":          kilometer,
	"kilometer":   kilometer,
	"kilometers":  kilometer,
	"mm":          millimeter,
	"millimeter":  millimeter,
	"millimeters": millimeter,
	"cm":          centimeter,
	"centimeter":  centimeter,
	"centimeters": centimeter,
	"in":          inch,
	"inch":        inch,
	"inches":      inch,
	"ft":          foot,
	"foot":        foot,
	"feet":        foot,
	"mi":          mile,
	"mile":        mile,
	"miles":       mile,
	"furlong":     furlong,
	"furlongs":    furlong,

	"g":         gram,
	"gram":      gram,
	"grams":     gram,
	"kg":        kilogram,
	"kilogram":  kilogram,
	"kilograms": kilogram,
	"lb":        pound,
	"lbs":       pound,
	"pound":     pound,
	"pounds":    pound,

	"c":          celcius,
	"celcius":    celcius,
	"f":          fahrenheit,
	"fahrenheit": fahrenheit,
	"kelvin":     kelvin,
	"k":          kelvin,
}

func parseUnit(s string) (unitType, error) {
	u, ok := unitMap[strings.ToLower(s)]
	if !ok {
		return nil, fmt.Errorf("Invalid unit %s", s)
	}
	return u, nil
}

func convert(num float64, from, to unitType) (float64, error) {
	switch from := from.(type) {
	case *tempertureUnit:
		if to, ok := to.(*tempertureUnit); ok {
			return from.convert(num, to), nil
		}
	case *weightUnit:
		if to, ok := to.(*weightUnit); ok {
			return from.convert(num, to), nil
		}
	case *lengthUnit:
		if to, ok := to.(*lengthUnit); ok {
			return from.convert(num, to), nil
		}
	}
	return 0, fmt.Errorf("Can't convert from %s to %s", from.name(), to.name())
}
