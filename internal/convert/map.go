package convert

var unitMap = map[string]UnitType{
	"m":           Meter,
	"meter":       Meter,
	"meters":      Meter,
	"km":          Kilometer,
	"kilometer":   Kilometer,
	"kilometers":  Kilometer,
	"mm":          Millimeter,
	"millimeter":  Millimeter,
	"millimeters": Millimeter,
	"cm":          Centimeter,
	"centimeter":  Centimeter,
	"centimeters": Centimeter,
	"nm":          Nanometer,
	"nanometer":   Nanometer,
	"nanometers":  Nanometer,
	"in":          Inch,
	"inch":        Inch,
	"inches":      Inch,
	"ft":          FootInch,
	"foot":        Foot,
	"feet":        Foot,
	"feet+inches": FootInch,
	"ftin":        FootInch,
	"ft+in":       FootInch,
	"yd":          Yard,
	"yard":        Yard,
	"yards":       Yard,
	"mi":          Mile,
	"mile":        Mile,
	"miles":       Mile,
	"furlong":     Furlong,
	"furlongs":    Furlong,
	"ly":          Lightyear,
	"lightyear":   Lightyear,
	"lightyears":  Lightyear,
	"g":           Gram,
	"gram":        Gram,
	"grams":       Gram,
	"kg":          Kilogram,
	"kilogram":    Kilogram,
	"kilograms":   Kilogram,
	"lb":          Pound,
	"lbs":         Pound,
	"pound":       Pound,
	"pounds":      Pound,
	"c":           Celsius,
	"celsius":     Celsius,
	"celcius":     Celsius,
	"f":           Fahrenheit,
	"fahrenheit":  Fahrenheit,
	"kelvin":      Kelvin,
	"k":           Kelvin,
	"mph":         MilesPerHour,
	"km/h":        KilometersPerHour,
	"kmph":        KilometersPerHour,
	"l":           Liter,
	"liter":       Liter,
	"liters":      Liter,
	"cl":          Centiliter,
	"centiliter":  Centiliter,
	"centiliters": Centiliter,
	"s":           Second,
	"sec":         Second,
	"secs":        Second,
	"second":      Second,
	"seconds":     Second,
	"min":         Minute,
	"mins":        Minute,
	"minute":      Minute,
	"minutes":     Minute,
	"hr":          Hour,
	"hrs":         Hour,
	"hour":        Hour,
	"hours":       Hour,
	"day":         Day,
	"days":        Day,
	"week":        Week,
	"weeks":       Week,
	"month":       Month,
	"months":      Month,
	"year":        Year,
	"years":       Year,
}
