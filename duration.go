package convert

import "github.com/martinlindhe/unit"

// DurationUnit is a unit of time
type DurationUnit = SimpleUnit[unit.Duration]

// Time units
var (
	Second = &DurationUnit{"s", from(unit.Second), unit.Duration.Seconds}
	Minute = &DurationUnit{"min", from(unit.Minute), unit.Duration.Minutes}
	Hour   = &DurationUnit{"hr", from(unit.Hour), unit.Duration.Hours}
	Day    = &DurationUnit{"days", from(unit.Day), unit.Duration.Days}
	Week   = &DurationUnit{"weeks", from(unit.Week), unit.Duration.Weeks}
	Month  = &DurationUnit{"months", from(unit.ThirtyDayMonth), unit.Duration.ThirtyDayMonths}
	Year   = &DurationUnit{"years", from(unit.JulianYear), unit.Duration.JulianYears}
)
