package main

import "github.com/martinlindhe/unit"

type tempertureUnit struct {
	unitCommon
	from func(float64) unit.Temperature
	to   func(unit.Temperature) float64
}

func (tu *tempertureUnit) fromFloat(f float64) unitVal {
	return temperatureVal{tu.from(f), tu}
}

var (
	celsius    = &tempertureUnit{"°C", unit.FromCelsius, unit.Temperature.Celsius}
	fahrenheit = &tempertureUnit{"°F", unit.FromFahrenheit, unit.Temperature.Fahrenheit}
	kelvin     = &tempertureUnit{"K", unit.FromKelvin, unit.Temperature.Kelvin}
)

type temperatureVal struct {
	v unit.Temperature
	u *tempertureUnit
}

func (tv temperatureVal) String() string {
	return simpleUnitString(tv.u.to(tv.v), tv.u)
}

func (tv temperatureVal) convert(to unitType) (unitVal, error) {
	if to, ok := to.(*tempertureUnit); ok {
		tv.u = to
		return tv, nil
	}
	return nil, convErr(tv.u, to)
}
