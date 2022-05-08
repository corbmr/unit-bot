package convert

var supportedUnits = map[UnitType][]string{
	// Length
	Meter:        {"m", "meter", "meters"},
	Kilometer:    {"km", "kilometer", "kilometers"},
	Millimeter:   {"mm", "millimeter", "millimeters"},
	Centimeter:   {"cm", "centimeter", "centimeters"},
	Nanometer:    {"nm", "nanometer", "nanometers"},
	Inch:         {"in", "inch", "inches"},
	FootInch:     {"ft", "feet+inches", "ftin", "ft+in"},
	Foot:         {"foot", "feet"},
	Yard:         {"yd", "yard", "yards"},
	Mile:         {"mi", "mile", "miles"},
	Furlong:      {"furlong", "furlongs"},
	Lightyear:    {"ly", "lightyear", "lightyears"},
	NauticalMile: {"nmi"},
	Fathom:       {"fathom", "fathoms"},

	// Mass
	Gram:     {"g", "gram", "grams"},
	Kilogram: {"kg", "kilogram", "kilograms"},
	Pound:    {"lb", "lbs", "pound", "pounds"},
	Stone:    {"st", "stone", "stones"},

	// Temperature
	Celsius:    {"c", "celcius", "celsius"},
	Fahrenheit: {"f", "fahrenheit"},
	Kelvin:     {"k", "kelvin"},

	// Speed
	MilesPerHour:      {"mph"},
	KilometersPerHour: {"kmh", "km/h", "kmph"},
	LightSpeed:        {"light", "lights", "lightspeed"},

	// Volume
	Liter:      {"l", "liter", "liters"},
	Centiliter: {"cl", "centiliter", "centiliters"},
	Milliliter: {"ml", "milliliter", "milliliters"},
	Gallon:     {"gal", "gals", "gallon", "gallons"},
	Quart:      {"qt", "quart", "quarts"},
	Pint:       {"pt", "pint", "pints"},
	Cup:        {"cup", "cups"},
	FlOunce:    {"oz", "floz", "ounce", "ounces"},
	Tablespoon: {"tbsp", "tablespoon", "tablespoons"},
	Teaspoon:   {"tsp", "teaspoon", "teaspoons"},

	// Duration
	Second: {"s", "sec", "secs", "second", "seconds"},
	Minute: {"min", "mins", "minute", "minutes"},
	Hour:   {"hr", "hrs", "hour", "hours"},
	Day:    {"day", "days"},
	Week:   {"wk", "week", "weeks"},
	Month:  {"month", "months"},
	Year:   {"yr", "year", "years"},
}
