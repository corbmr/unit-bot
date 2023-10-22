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

type LengthVal struct {
	SimpleUnitValue[unit.Length]
}

// Length units
var (
	Meter        = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"m", from(unit.Meter), unit.Length.Meters}}
	Kilometer    = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"km", from(unit.Kilometer), unit.Length.Kilometers}}
	Millimeter   = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"mm", from(unit.Millimeter), unit.Length.Millimeters}}
	Centimeter   = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"cm", from(unit.Centimeter), unit.Length.Centimeters}}
	Nanometer    = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"nm", from(unit.Nanometer), unit.Length.Nanometers}}
	Inch         = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"in", from(unit.Inch), unit.Length.Inches}}
	Foot         = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"ft", from(unit.Foot), unit.Length.Feet}}
	Yard         = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"yd", from(unit.Yard), unit.Length.Yards}}
	Mile         = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"miles", from(unit.Mile), unit.Length.Miles}}
	Furlong      = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"furlongs", from(unit.Furlong), unit.Length.Furlongs}}
	Lightyear    = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"ly", from(unit.LightYear), unit.Length.LightYears}}
	NauticalMile = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"nautical mile", from(unit.NauticalMile), unit.Length.NauticalMiles}}
	Fathom       = &LengthUnit{SimpleUnit: SimpleUnit[unit.Length]{"fathoms", from(unit.Fathom), unit.Length.Fathoms}}

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
