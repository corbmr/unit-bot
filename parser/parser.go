package parser

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var wsPattern = regexp.MustCompile(`^\s*`)

func ws(s []byte) int {
	return len(wsPattern.Find(s))
}

// Res is the result of parsing
// The zero value for Res signals that the parsing was unsucessful
type Res struct {
	// Value return from the parser
	V interface{}
	// Number of bytes consumed
	N int
	// Whether the parsing was successful
	// If Ok is false, V should be nil and N should be 0
	Ok bool
}

// Parser scans bytes and gives back a result
type Parser func([]byte) Res

// Token scans for a pattern, skipping leading whitespace if necessary
func Token(pattern string) Parser {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) Res {
		n := ws(s)
		match := regex.Find(s[n:])
		if match == nil {
			return Res{}
		}

		return Res{string(match), n + len(match), true}
	}
}

// TokenE is like Token but without skipping leading whitespace
func TokenE(pattern string) Parser {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) Res {
		match := regex.Find(s)
		if match == nil {
			return Res{}
		}

		return Res{string(match), len(match), true}
	}
}

// Sub matches sub expressions in a regex pattern
func Sub(pattern string) Parser {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) Res {
		n := ws(s)
		match := regex.FindSubmatch(s[n:])
		if match == nil {
			return Res{}
		}
		matches := make(map[string]string)
		for i, n := range regex.SubexpNames() {
			if len(n) > 0 {
				matches[n] = string(match[i])
			}
		}
		return Res{matches, n + len(match[0]), true}
	}
}

// All matches all parsers in order
// Will result in a slice of all of the parsed values
func All(ps ...Parser) Parser {
	return func(s []byte) Res {
		var (
			vs  []interface{}
			sum int
		)
		for _, p := range ps {
			res := p(s[sum:])
			if !res.Ok {
				return Res{}
			}
			vs = append(vs, res.V)
			sum += res.N
		}
		return Res{vs, sum, true}
	}
}

// Any matches any of the parsers, tested in order
func Any(ps ...Parser) Parser {
	return func(s []byte) Res {
		for _, p := range ps {
			res := p(s)
			if res.Ok {
				return res
			}
		}
		return Res{}
	}
}

// Atom scans for a single atom
func Atom(val string) Parser {
	b := []byte(val)
	return func(s []byte) Res {
		if n := ws(s); bytes.HasPrefix(s[n:], b) {
			return Res{nil, n + len(b), true}
		}
		return Res{}
	}
}

// AtomE is like Atom but without skipping leading whitespace
func AtomE(val string) Parser {
	b := []byte(val)
	return func(s []byte) Res {
		if bytes.HasPrefix(s, b) {
			return Res{nil, len(b), true}
		}
		return Res{}
	}
}

// RuneIn matches a single rune of the ones in val
func RuneIn(val string) Parser {
	return func(s []byte) Res {
		w := ws(s)
		r, n := utf8.DecodeRune(s[w:])
		if strings.ContainsRune(val, r) {
			return Res{r, w + n, true}
		}
		return Res{}
	}
}

// None is type representing missing from Opt
type None struct{}

// Opt turns a parser into an optional parser
// Will return a result that's either Missing or the result itself
func (p Parser) Opt() Parser {
	return func(s []byte) Res {
		res := p(s)
		if !res.Ok {
			return Res{None{}, 0, true}
		}
		return res
	}
}

// Or is like Opt but returns the value given if the parser is unsuccessful
func (p Parser) Or(or interface{}) Parser {
	return func(s []byte) Res {
		res := p(s)
		if !res.Ok {
			return Res{or, 0, true}
		}
		return res
	}
}

// Map maps the result of a parser to a different result
// If the Mapper returns nil, the parser returns as invalid
func (p Parser) Map(f Mapper) Parser {
	return func(s []byte) Res {
		if res := p(s); res.Ok {
			if v := f(res.V); v != nil {
				return Res{v, res.N, true}
			}
		}
		return Res{}
	}
}

// Float is a float parser
var Float = Token(`[+-]?\d+([.]\d*)?([eE][+-]?\d+)?`).Map(MapFloat)

// Int is an integer parser
var Int = Token(`[+-]?\d+`).Map(MapInt)

// Mapper is a function for mapping parser results
type Mapper = func(interface{}) interface{}

// Index creates a mapper that maps the result to an index in a slice
func Index(i int) Mapper {
	return func(v interface{}) interface{} {
		return v.([]interface{})[i]
	}
}

// MapFloat maps the string result to a float
func MapFloat(v interface{}) interface{} {
	f, err := strconv.ParseFloat(v.(string), 64)
	if err != nil {
		panic(err)
	}
	return f
}

// MapInt maps string result to an int
func MapInt(v interface{}) interface{} {
	i, err := strconv.Atoi(v.(string))
	if err != nil {
		panic(err)
	}
	return i
}
