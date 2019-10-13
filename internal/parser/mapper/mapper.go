package mapper

import "strconv"

// Mapper is a function for mapping parser results
type Mapper func(interface{}) interface{}

// Index creates a mapper that maps the result to an index in a slice
func Index(i int) Mapper {
	return func(v interface{}) interface{} {
		return v.([]interface{})[i]
	}
}

// Float maps the string result to a float
func Float(v interface{}) interface{} {
	f, err := strconv.ParseFloat(v.(string), 64)
	if err != nil {
		panic(err)
	}
	return f
}

// Int maps string result to an int
func Int(v interface{}) interface{} {
	i, err := strconv.Atoi(v.(string))
	if err != nil {
		panic(err)
	}
	return i
}
