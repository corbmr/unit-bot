package convert

import "github.com/martinlindhe/unit"

// SpeedUnit is a unit of speed
type SpeedUnit struct {
	unitCommon
	speed   unit.Speed
	toFloat func(unit.Speed) float64
}

// FromFloat implements SimpleUnit
func (su *SpeedUnit) FromFloat(f float64) UnitVal {
	return SpeedVal{unit.Speed(f) * su.speed, su}
}

// Speed units
var (
	MilesPerHour      = &SpeedUnit{"mph", unit.MilesPerHour, unit.Speed.MilesPerHour}
	KilometersPerHour = &SpeedUnit{"km/h", unit.KilometersPerHour, unit.Speed.KilometersPerHour}
	LightSpeed        = &SpeedUnit{"C", unit.SpeedOfLight, unit.Speed.SpeedOfLight}
)

// SpeedVal is a speed value with unit
type SpeedVal struct {
	V unit.Speed
	U *SpeedUnit
}

func (sv SpeedVal) String() string {
	return simpleUnitString(sv.U.toFloat(sv.V), sv.U)
}

// Convert implements UnitVal conversion
func (sv SpeedVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*SpeedUnit); ok {
		sv.U = to
		return sv, nil
	}
	return nil, ErrorConversion{sv.U, to}
}
