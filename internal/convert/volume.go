package convert

import "github.com/martinlindhe/unit"

// VolumeUnit is a volumetric unit
type VolumeUnit struct {
	unitCommon
	volume unit.Volume
	to     func(unit.Volume) float64
}

// FromFloat implements SimpleUnit
func (vu *VolumeUnit) FromFloat(f float64) UnitVal {
	return VolumeVal{unit.Volume(f) * vu.volume, vu}
}

// Volumetric units
var (
	Milliliter = &VolumeUnit{"ml", unit.Milliliter, unit.Volume.Milliliters}
	Centiliter = &VolumeUnit{"cl", unit.Centiliter, unit.Volume.Centiliters}
	Liter      = &VolumeUnit{"l", unit.Liter, unit.Volume.Liters}

	Gallon     = &VolumeUnit{"gal", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
	Quart      = &VolumeUnit{"quart", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
	Pint       = &VolumeUnit{"pint", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
	Cup        = &VolumeUnit{"cup", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
	Ounce      = &VolumeUnit{"oz", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
	TableSpoon = &VolumeUnit{"tbsp", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
	TeaSpoon   = &VolumeUnit{"tsp", unit.USLiquidGallon, unit.Volume.USLiquidGallons}
)

// VolumeVal is a volumetric value with unit
type VolumeVal struct {
	V unit.Volume
	U *VolumeUnit
}

func (vv VolumeVal) String() string {
	return simpleUnitString(vv.U.to(vv.V), vv.U)
}

// Convert implements UnitVal conversion
func (vv VolumeVal) Convert(to UnitType) (UnitVal, error) {
	if to, ok := to.(*VolumeUnit); ok {
		vv.U = to
		return vv, nil
	}
	return nil, ErrorConversion{vv.U, to}
}
