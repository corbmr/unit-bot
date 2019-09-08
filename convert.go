package main

import (
	"fmt"
	"strings"

	"github.com/martinlindhe/unit"
)

type unitType int

const (
	length unitType = iota
	weight
	speed
	temperature
)

// Unit is represents a unit
type Unit struct {
	utype unitType
	name  string
}

var (
	meter      = &Unit{length, "m"}
	kilometer  = &Unit{length, "km"}
	millimeter = &Unit{length, "mm"}
	centimeter = &Unit{length, "cm"}
	inch       = &Unit{length, "in"}
	foot       = &Unit{length, "ft"}
	mile       = &Unit{length, "miles"}
	furlong    = &Unit{length, "furlongs"}

	gram     = &Unit{weight, "g"}
	kilogram = &Unit{weight, "kg"}
	pound    = &Unit{weight, "lbs"}

	celcius    = &Unit{temperature, "°C"}
	fahrenheit = &Unit{temperature, "°F"}
)

var unitMap = map[string]*Unit{
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
	"mile":        mile,
	"miles":       mile,
	"g":           gram,
	"gram":        gram,
	"grams":       gram,
	"kg":          kilogram,
	"kilogram":    kilogram,
	"kilograms":   kilogram,
	"lb":          pound,
	"lbs":         pound,
	"pound":       pound,
	"pounds":      pound,
	"c":           celcius,
	"celcius":     celcius,
	"f":           fahrenheit,
	"fahrenheit":  fahrenheit,
	"furlong":     furlong,
	"furlongs":    furlong,
}

func parseUnit(s string) (*Unit, error) {
	u, ok := unitMap[strings.ToLower(s)]
	if !ok {
		return nil, fmt.Errorf("Invalid unit %s", s)
	}
	return u, nil
}

func convert(num float64, unitFrom, unitTo *Unit) (float64, error) {
	if unitFrom.utype != unitTo.utype {
		return 0, fmt.Errorf("Can't convert from %s to %s", unitFrom.name, unitTo.name)
	}

	switch unitFrom.utype {
	case length:
		return convertLength(num, unitFrom, unitTo), nil
	case weight:
		return convertWeight(num, unitFrom, unitTo), nil
	case temperature:
		return convertTemperature(num, unitFrom, unitTo), nil
	default:
		return 0, nil
	}
}

func convertLength(num float64, unitFrom, unitTo *Unit) float64 {
	length := unit.Length(num)

	var from unit.Length
	switch unitFrom {
	case meter:
		from = length * unit.Meter
	case kilometer:
		from = length * unit.Kilometer
	case millimeter:
		from = length * unit.Millimeter
	case centimeter:
		from = length * unit.Centimeter
	case foot:
		from = length * unit.Foot
	case inch:
		from = length * unit.Inch
	case mile:
		from = length * unit.Mile
	case furlong:
		from = length * unit.Furlong
	}

	var to float64
	switch unitTo {
	case meter:
		to = from.Meters()
	case kilometer:
		to = from.Kilometers()
	case millimeter:
		to = from.Millimeters()
	case centimeter:
		to = from.Centimeters()
	case foot:
		to = from.Feet()
	case inch:
		to = from.Inches()
	case mile:
		to = from.Miles()
	case furlong:
		to = from.Furlongs()
	}
	return to
}

func convertWeight(num float64, unitFrom, unitTo *Unit) float64 {
	weight := unit.Mass(num)

	var from unit.Mass
	switch unitFrom {
	case gram:
		from = weight * unit.Gram
	case kilogram:
		from = weight * unit.Kilogram
	case pound:
		from = weight * unit.AvoirdupoisPound
	}

	var to float64
	switch unitTo {
	case gram:
		to = from.Grams()
	case kilogram:
		to = from.Kilograms()
	case pound:
		to = from.AvoirdupoisPounds()
	}
	return to
}

func convertTemperature(num float64, unitFrom, unitTo *Unit) float64 {
	var from unit.Temperature
	switch unitFrom {
	case celcius:
		from = unit.FromCelsius(num)
	case fahrenheit:
		from = unit.FromFahrenheit(num)
	}

	var to float64
	switch unitTo {
	case celcius:
		to = from.Celsius()
	case fahrenheit:
		to = from.Fahrenheit()
	}
	return to
}
