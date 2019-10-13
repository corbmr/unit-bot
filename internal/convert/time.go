package convert

import "github.com/martinlindhe/unit"

// TimeUnit is a unit of time
type TimeUnit struct {
	unitCommon
	time unit.Duration
	to   func(unit.Duration) float64
}

// FromFloat implements SimpleUnit
func (tu *TimeUnit) FromFloat(f float64) UnitVal {
	return TimeVal{unit.Duration(f) * tu.time, tu}
}

// Time units
var (
	Second = &TimeUnit{"s", unit.Second, unit.Duration.Seconds}
	Hour   = &TimeUnit{"hr", unit.Hour, unit.Duration.Hours}
	Day    = &TimeUnit{"days", unit.Day, unit.Duration.Days}
	Week   = &TimeUnit{"weeks", unit.Week, unit.Duration.Weeks}
	Month  = &TimeUnit{"months", unit.ThirtyDayMonth, unit.Duration.ThirtyDayMonths}
	Year   = &TimeUnit{"years", unit.JulianYear, unit.Duration.JulianYears}
)

// TimeVal is a time value with unit
type TimeVal struct {
	V unit.Duration
	U *TimeUnit
}

func (tv TimeVal) String() string {
	return simpleUnitString(tv.U.to(tv.V), tv.U)
}

// Convert implements UnitVal conversion
func (tv TimeVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*TimeUnit); ok {
		tv.U = to
		return tv, nil
	}
	return nil, convErr(tv.U, to)
}
