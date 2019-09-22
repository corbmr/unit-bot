package main

import "github.com/martinlindhe/unit"

type speedUnit struct {
	unitCommon
	speed unit.Speed
	to    func(unit.Speed) float64
}

func (su *speedUnit) fromFloat(f float64) unitVal {
	return speedVal{unit.Speed(f) * su.speed, su}
}

var (
	milesPerHour      = &speedUnit{"mph", unit.MilesPerHour, unit.Speed.MilesPerHour}
	kilometersPerHour = &speedUnit{"km/h", unit.KilometersPerHour, unit.Speed.KilometersPerHour}
)

type speedVal struct {
	v unit.Speed
	u *speedUnit
}

func (sv speedVal) String() string {
	return simpleUnitString(sv.u.to(sv.v), sv.u)
}

func (sv speedVal) convert(to unitType) (unitVal, error) {
	if to, ok := to.(*speedUnit); ok {
		sv.u = to
		return sv, nil
	}
	return nil, convErr(sv.u, to)
}
