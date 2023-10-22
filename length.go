package convert

import (
	"fmt"
	"math"

	"github.com/martinlindhe/unit"
)

// LengthUnit is a unit of length
type LengthUnit struct {
	SimpleUnit[unit.Length]
}

type LengthVal struct{ SimpleUnitValue[unit.Length] }

// Length units
var (
	Meter        = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "m", from(unit.Meter), unit.Length.Meters}}
	Kilometer    = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "km", from(unit.Kilometer), unit.Length.Kilometers}}
	Millimeter   = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "mm", from(unit.Millimeter), unit.Length.Millimeters}}
	Centimeter   = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "cm", from(unit.Centimeter), unit.Length.Centimeters}}
	Nanometer    = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "nm", from(unit.Nanometer), unit.Length.Nanometers}}
	Inch         = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "in", from(unit.Inch), unit.Length.Inches}}
	Foot         = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "ft", from(unit.Foot), unit.Length.Feet}}
	Yard         = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "yd", from(unit.Yard), unit.Length.Yards}}
	Mile         = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "miles", from(unit.Mile), unit.Length.Miles}}
	Furlong      = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "furlongs", from(unit.Furlong), unit.Length.Furlongs}}
	Lightyear    = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "ly", from(unit.LightYear), unit.Length.LightYears}}
	NauticalMile = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "nautical mile", from(unit.NauticalMile), unit.Length.NauticalMiles}}
	Fathom       = &LengthUnit{SimpleUnit[unit.Length]{UnitDimensionLength, "fathoms", from(unit.Fathom), unit.Length.Fathoms}}

	FootInch = &FootInchUnit{}
)

func (u *LengthUnit) FromFloat(f float64) UnitVal {
	return LengthVal{u.SimpleUnit.FromFloat(f).(SimpleUnitValue[unit.Length])}
}

// Convert implements UnitVal conversion
func (lv LengthVal) Convert(to UnitType) (UnitVal, error) {
	switch to := to.(type) {
	case *LengthUnit:
		lv.unit = &to.SimpleUnit
		return lv, nil
	case *FootInchUnit:
		feet, fraction := math.Modf(lv.value.Feet())
		inches := (unit.Length(fraction) * unit.Foot).Inches()
		return FootInchVal{feet, inches}, nil
	default:
		return nil, ErrorConversion{lv.unit, to}
	}
}

// FootInchUnit is a unit of both feet + inches
type FootInchUnit struct{}

// FootInchVal is a UnitVal of both feet + inches
type FootInchVal struct {
	Feet, Inches float64
}

// FromFloat implements UnitType
func (FootInchUnit) FromFloat(f float64) UnitVal {
	return FootInchVal{Feet: f, Inches: 0}
}

func (FootInchUnit) String() string {
	return "feet + inches"
}

func (FootInchUnit) Dimension() UnitDimension {
	return UnitDimensionLength
}

func (val FootInchVal) String() string {
	if val.Inches == 0 {
		return Foot.FromFloat(val.Feet).String()
	}
	return fmt.Sprintf(`%.f' %.f"`, val.Feet, val.Inches)
}

// Convert implements UnitVal conversion
func (val FootInchVal) Convert(to UnitType) (UnitVal, error) {
	switch to := to.(type) {
	case *FootInchUnit:
		return val, nil
	case *LengthUnit:
		feet := unit.Length(val.Feet) * unit.Foot
		inches := unit.Length(val.Inches) * unit.Inch
		return LengthVal{SimpleUnitValue[unit.Length]{feet + inches, &to.SimpleUnit}}, nil
	default:
		return nil, ErrorConversion{FootInch, to}
	}
}

func (FootInchVal) Unit() UnitType {
	return FootInch
}
