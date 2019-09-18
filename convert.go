package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type unitType interface {
	name() string
}

type unitCommon string

func (c unitCommon) name() string {
	return string(c)
}

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
	"nm":          nanometer,
	"nanometer":   nanometer,
	"nanometers":  nanometer,
	"in":          inch,
	"inch":        inch,
	"inches":      inch,
	"ft":          foot,
	"foot":        foot,
	"feet":        foot,
	"yd":          yard,
	"yard":        yard,
	"yards":       yard,
	"mi":          mile,
	"mile":        mile,
	"miles":       mile,
	"furlong":     furlong,
	"furlongs":    furlong,
	"ly":          lightyear,
	"lightyear":   lightyear,
	"lightyears":  lightyear,
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
	"c":           celsius,
	"celsius":     celsius,
	"celcius":     celsius,
	"f":           fahrenheit,
	"fahrenheit":  fahrenheit,
	"kelvin":      kelvin,
	"k":           kelvin,
}

const usage = `Usage: !conv [from][unit] to [unit]`

// TODO: write a real parser, this is gross
var inputRegex = regexp.MustCompile(
	`^(((?P<num>[+-]?\d+([.]\d*)?([eE][+-]?\d+)?)\s*(?P<from>\w+))|((?P<feet>[+-]?\d+)'\s*(?P<inches>\d+"?)?))\s+to\s+(?P<to>\w+)`)

func generateResponse(inp string) (string, error) {
	match := findNamed(inputRegex, inp)
	if match == nil {
		return "", fmt.Errorf(usage)
	}

	to, err := parseUnit(match["to"])
	if err != nil {
		return "", err
	}

	// Normal request branch
	if match["num"] != "" && match["from"] != "" {
		num, _ := strconv.ParseFloat(match["num"], 64)

		from, err := parseUnit(match["from"])
		if err != nil {
			return "", err
		}

		conv, err := convert(num, from, to)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s %s = %.6g %s",
			match["num"], from.name(), conv, to.name()), nil
	}

	// Branch using X'Y" syntax for ft+in
	if feet := match["feet"]; feet != "" {

		toLength, ok := to.(*lengthUnit)
		if !ok {
			return "", fmt.Errorf("Can't convert ft to %s", to.name())
		}

		var f, i float64
		f, _ = strconv.ParseFloat(feet, 64)
		if inches := match["inches"]; inches != "" {
			i, _ = strconv.ParseFloat(inches, 64)
		}

		conv := convFeetInches(f, i, toLength)
		return fmt.Sprintf(`%.0f'%.0f" = %.6g %s`,
			f, i, conv, toLength.name()), nil
	}

	// This should never happen if the regex is correct
	panic("Unexpected state")
}

func findNamed(r *regexp.Regexp, s string) map[string]string {
	match := r.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	names := r.SubexpNames()
	out := make(map[string]string)
	for i, n := range names {
		if n != "" {
			out[n] = match[i]
		}
	}

	return out
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
	case *massUnit:
		if to, ok := to.(*massUnit); ok {
			return from.convert(num, to), nil
		}
	case *lengthUnit:
		if to, ok := to.(*lengthUnit); ok {
			return from.convert(num, to), nil
		}
	}
	return 0, fmt.Errorf("Can't convert from %s to %s", from.name(), to.name())
}
