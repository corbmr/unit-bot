package convert

import "github.com/martinlindhe/unit"

// SpeedUnit is a unit of speed
type SpeedUnit = SimpleUnit[unit.Speed]

// Speed units
var (
	MilesPerHour      = &SpeedUnit{UnitDimensionSpeed, "mph", from(unit.MilesPerHour), unit.Speed.MilesPerHour}
	KilometersPerHour = &SpeedUnit{UnitDimensionSpeed, "km/h", from(unit.KilometersPerHour), unit.Speed.KilometersPerHour}
	LightSpeed        = &SpeedUnit{UnitDimensionSpeed, "C", from(unit.SpeedOfLight), unit.Speed.SpeedOfLight}
)
