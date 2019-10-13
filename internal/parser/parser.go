package parser

import (
	"bytes"
	"regexp"

	"github.com/corbmr/unit-bot/internal/parser/mapper"
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
		return Res{res.V, res.N, true}
	}
}

// Map maps the result of a parser to a different result
func (p Parser) Map(f mapper.Mapper) Parser {
	return func(s []byte) Res {
		if res := p(s); res.Ok {
			return Res{f(res.V), res.N, true}
		}
		return Res{}
	}
}

// Float is a float parser
var Float = Token(`[+-]?\d+([.]\d*)?([eE][+-]?\d+)?`).Map(mapper.Float)

// Int is an integer parser
var Int = Token(`[+-]?\d+`).Map(mapper.Int)
