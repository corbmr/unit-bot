package convert

import "github.com/martinlindhe/unit"

// MassUnit is a unit of mass
type MassUnit struct {
	unitCommon
	mass    unit.Mass
	toFloat func(unit.Mass) float64
}

// FromFloat implements SimpleUnit
func (mu *MassUnit) FromFloat(f float64) UnitVal {
	return MassVal{unit.Mass(f) * mu.mass, mu}
}

// Mass units
var (
	Gram     = &MassUnit{"g", unit.Gram, unit.Mass.Grams}
	Kilogram = &MassUnit{"kg", unit.Kilogram, unit.Mass.Kilograms}
	Pound    = &MassUnit{"lbs", unit.AvoirdupoisPound, unit.Mass.AvoirdupoisPounds}
	Stone    = &MassUnit{"stones", unit.UkStone, unit.Mass.UkStones}
)

// MassVal is a mass value with unit
type MassVal struct {
	V unit.Mass
	U *MassUnit
}

func (mv MassVal) String() string {
	return simpleUnitString(mv.U.toFloat(mv.V), mv.U)
}

// Convert implements UnitVal conversion
func (mv MassVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*MassUnit); ok {
		mv.U = to
		return mv, nil
	}
	return nil, ErrorConversion{mv.U, to}
}
