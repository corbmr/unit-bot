package convert

import "github.com/martinlindhe/unit"

// MassUnit is a unit of mass
type MassUnit = SimpleUnit[unit.Mass]

// Mass units
var (
	Gram     = &MassUnit{UnitDimensionMass, "g", from(unit.Gram), unit.Mass.Grams}
	Kilogram = &MassUnit{UnitDimensionMass, "kg", from(unit.Kilogram), unit.Mass.Kilograms}
	Pound    = &MassUnit{UnitDimensionMass, "lbs", from(unit.AvoirdupoisPound), unit.Mass.AvoirdupoisPounds}
	Stone    = &MassUnit{UnitDimensionMass, "stones", from(unit.UkStone), unit.Mass.UkStones}
)
