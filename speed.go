package convert

import "github.com/martinlindhe/unit"

// SpeedUnit is a unit of speed
type SpeedUnit = SimpleUnit[unit.Speed]

// Speed units
var (
	MilesPerHour      = &SpeedUnit{"mph", from(unit.MilesPerHour), unit.Speed.MilesPerHour}
	KilometersPerHour = &SpeedUnit{"km/h", from(unit.KilometersPerHour), unit.Speed.KilometersPerHour}
	LightSpeed        = &SpeedUnit{"C", from(unit.SpeedOfLight), unit.Speed.SpeedOfLight}
)
