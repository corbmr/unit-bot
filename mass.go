package main

import "github.com/martinlindhe/unit"

type massUnit struct {
	unitCommon
	mass unit.Mass
	to   func(unit.Mass) float64
}

func (mu *massUnit) fromFloat(f float64) unitVal {
	return massVal{unit.Mass(f) * mu.mass, mu}
}

var (
	gram     = &massUnit{"g", unit.Gram, unit.Mass.Grams}
	kilogram = &massUnit{"kg", unit.Kilogram, unit.Mass.Kilograms}
	pound    = &massUnit{"lbs", unit.AvoirdupoisPound, unit.Mass.AvoirdupoisPounds}
)

type massVal struct {
	v unit.Mass
	u *massUnit
}

func (mv massVal) String() string {
	return simpleUnitString(mv.u.to(mv.v), mv.u)
}

func (mv massVal) convert(to unitType) (unitVal, error) {
	if to, ok := to.(*massUnit); ok {
		mv.u = to
		return mv, nil
	}
	return nil, convErr(mv.u, to)
}
