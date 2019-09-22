package main

import "github.com/martinlindhe/unit"

type volumeUnit struct {
	unitCommon
	volume unit.Volume
	to     func(unit.Volume) float64
}

func (vu *volumeUnit) fromFloat(f float64) unitVal {
	return volumeVal{unit.Volume(f) * vu.volume, vu}
}

var (
	centiliter = &volumeUnit{"cl", unit.Centiliter, unit.Volume.Centiliters}
	liter      = &volumeUnit{"l", unit.Liter, unit.Volume.Liters}
)

type volumeVal struct {
	v unit.Volume
	u *volumeUnit
}

func (vv volumeVal) String() string {
	return simpleUnitString(vv.u.to(vv.v), vv.u)
}

func (vv volumeVal) convert(to unitType) (unitVal, error) {
	if to, ok := to.(*volumeUnit); ok {
		vv.u = to
		return vv, nil
	}
	return nil, convErr(vv.u, to)
}
