package convert

import "github.com/martinlindhe/unit"

// TemperatureUnit is a unit of temperture
type TemperatureUnit = SimpleUnit[unit.Temperature]

// Temperature units
var (
	Celsius    = &TemperatureUnit{UnitDimensionTemperature, "°C", unit.FromCelsius, unit.Temperature.Celsius}
	Fahrenheit = &TemperatureUnit{UnitDimensionTemperature, "°F", unit.FromFahrenheit, unit.Temperature.Fahrenheit}
	Kelvin     = &TemperatureUnit{UnitDimensionTemperature, "K", unit.FromKelvin, unit.Temperature.Kelvin}
)
