package main

import "github.com/martinlindhe/unit"

type tempertureUnit struct {
	unitCommon
	from func(float64) unit.Temperature
	to   func(unit.Temperature) float64
}

func (from *tempertureUnit) convert(f float64, to *tempertureUnit) float64 {
	return to.to(from.from(f))
}

var (
	celsius    = &tempertureUnit{"°C", unit.FromCelsius, unit.Temperature.Celsius}
	fahrenheit = &tempertureUnit{"°F", unit.FromFahrenheit, unit.Temperature.Fahrenheit}
	kelvin     = &tempertureUnit{"K", unit.FromKelvin, unit.Temperature.Kelvin}
)
