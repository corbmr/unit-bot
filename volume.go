package convert

import "github.com/martinlindhe/unit"

// VolumeUnit is a volumetric unit
type VolumeUnit = SimpleUnit[unit.Volume]

// Volumetric units
var (
	Milliliter = &VolumeUnit{UnitDimensionVolume, "ml", from(unit.Milliliter), unit.Volume.Milliliters}
	Centiliter = &VolumeUnit{UnitDimensionVolume, "cl", from(unit.Centiliter), unit.Volume.Centiliters}
	Liter      = &VolumeUnit{UnitDimensionVolume, "l", from(unit.Liter), unit.Volume.Liters}

	Gallon     = &VolumeUnit{UnitDimensionVolume, "gal", from(unit.USLiquidGallon), unit.Volume.USLiquidGallons}
	Quart      = &VolumeUnit{UnitDimensionVolume, "quart", from(unit.USLiquidQuart), unit.Volume.USLiquidQuarts}
	Pint       = &VolumeUnit{UnitDimensionVolume, "pint", from(unit.USLiquidPint), unit.Volume.USLiquidPints}
	Cup        = &VolumeUnit{UnitDimensionVolume, "cup", from(unit.USCup), unit.Volume.USCups}
	FlOunce    = &VolumeUnit{UnitDimensionVolume, "fl oz", from(unit.USFluidOunce), unit.Volume.USFluidOunces}
	Tablespoon = &VolumeUnit{UnitDimensionVolume, "tbsp", from(unit.USTableSpoon), unit.Volume.USTableSpoons}
	Teaspoon   = &VolumeUnit{UnitDimensionVolume, "tsp", from(unit.USTeaSpoon), unit.Volume.USTeaSpoons}
)
