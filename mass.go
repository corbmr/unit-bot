package main

import "github.com/martinlindhe/unit"

type massUnit struct {
	unitCommon
	weight unit.Mass
	to     func(unit.Mass) float64
}

func (from *massUnit) convert(f float64, to *massUnit) float64 {
	return to.to(unit.Mass(f) * from.weight)
}

var (
	gram     = &massUnit{"g", unit.Gram, unit.Mass.Grams}
	kilogram = &massUnit{"kg", unit.Kilogram, unit.Mass.Kilograms}
	pound    = &massUnit{"lbs", unit.AvoirdupoisPound, unit.Mass.AvoirdupoisPounds}
)
