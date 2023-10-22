package convert

import "github.com/martinlindhe/unit"

// DurationUnit is a unit of time
type DurationUnit = SimpleUnit[unit.Duration]

// Time units
var (
	Second = &DurationUnit{UnitDimensionDuration, "s", from(unit.Second), unit.Duration.Seconds}
	Minute = &DurationUnit{UnitDimensionDuration, "min", from(unit.Minute), unit.Duration.Minutes}
	Hour   = &DurationUnit{UnitDimensionDuration, "hr", from(unit.Hour), unit.Duration.Hours}
	Day    = &DurationUnit{UnitDimensionDuration, "days", from(unit.Day), unit.Duration.Days}
	Week   = &DurationUnit{UnitDimensionDuration, "weeks", from(unit.Week), unit.Duration.Weeks}
	Month  = &DurationUnit{UnitDimensionDuration, "months", from(unit.ThirtyDayMonth), unit.Duration.ThirtyDayMonths}
	Year   = &DurationUnit{UnitDimensionDuration, "years", from(unit.JulianYear), unit.Duration.JulianYears}
)
