package convert

import "github.com/martinlindhe/unit"

// VolumeUnit is a volumetric unit
type VolumeUnit = SimpleUnit[unit.Volume]

// Volumetric units
var (
	Milliliter = &VolumeUnit{"ml", from(unit.Milliliter), unit.Volume.Milliliters}
	Centiliter = &VolumeUnit{"cl", from(unit.Centiliter), unit.Volume.Centiliters}
	Liter      = &VolumeUnit{"l", from(unit.Liter), unit.Volume.Liters}

	Gallon     = &VolumeUnit{"gal", from(unit.USLiquidGallon), unit.Volume.USLiquidGallons}
	Quart      = &VolumeUnit{"quart", from(unit.USLiquidQuart), unit.Volume.USLiquidQuarts}
	Pint       = &VolumeUnit{"pint", from(unit.USLiquidPint), unit.Volume.USLiquidPints}
	Cup        = &VolumeUnit{"cup", from(unit.USCup), unit.Volume.USCups}
	FlOunce    = &VolumeUnit{"fl oz", from(unit.USFluidOunce), unit.Volume.USFluidOunces}
	Tablespoon = &VolumeUnit{"tbsp", from(unit.USTableSpoon), unit.Volume.USTableSpoons}
	Teaspoon   = &VolumeUnit{"tsp", from(unit.USTeaSpoon), unit.Volume.USTeaSpoons}
)
