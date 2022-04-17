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

// Parser scans bytes and returns result a result, number of bytes consumed, and whether the parse was successful
type Parser[T any] func([]byte) (T, int, bool)

// Token scans for a pattern, skipping leading whitespace if necessary
func Token(pattern string) Parser[string] {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) (string, int, bool) {
		n := ws(s)
		match := regex.Find(s[n:])
		if match == nil {
			return "", 0, false
		}

		return string(match), n + len(match), true
	}
}

// TokenE is like Token but without skipping leading whitespace
func TokenE(pattern string) Parser[string] {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) (string, int, bool) {
		match := regex.Find(s)
		if match == nil {
			return "", 0, false
		}

		return string(match), len(match), true
	}
}

// Sub matches sub expressions in a regex pattern
func Sub(pattern string) Parser[map[string]string] {
	if pattern[0] != '^' {
		pattern = "^" + pattern
	}
	regex := regexp.MustCompile(pattern)
	return func(s []byte) (map[string]string, int, bool) {
		n := ws(s)
		match := regex.FindSubmatch(s[n:])
		if match == nil {
			return nil, 0, false
		}
		matches := make(map[string]string)
		for i, n := range regex.SubexpNames() {
			if len(n) > 0 {
				matches[n] = string(match[i])
			}
		}
		return matches, n + len(match[0]), true
	}
}

// All matches all parsers in order
// Will result in a slice of all of the parsed values
func All[T any](ps ...Parser[T]) Parser[[]T] {
	return func(s []byte) ([]T, int, bool) {
		var (
			vs  []T
			sum int
		)
		for _, p := range ps {
			res, n, ok := p(s[sum:])
			if !ok {
				return nil, 0, false
			}
			vs = append(vs, res)
			sum += n
		}
		return vs, sum, true
	}
}

func Parse2[A, B, C any](pa Parser[A], pb Parser[B], f func(A, B) C) Parser[C] {
	return func(s []byte) (C, int, bool) {
		var c C
		a, n1, ok := pa(s)
		if !ok {
			return c, 0, false
		}
		b, n2, ok := pb(s[n1:])
		if !ok {
			return c, 0, false
		}
		return f(a, b), n1 + n2, true
	}
}

func Parse3[A, B, C, D any](pa Parser[A], pb Parser[B], pc Parser[C], f func(A, B, C) D) Parser[D] {
	return func(s []byte) (D, int, bool) {
		var d D
		a, n1, ok := pa(s)
		if !ok {
			return d, 0, false
		}
		b, n2, ok := pb(s[n1:])
		if !ok {
			return d, 0, false
		}
		c, n3, ok := pc(s[n1+n2:])
		if !ok {
			return d, 0, false
		}
		return f(a, b, c), n1 + n2 + n3, true
	}
}

// First matches the first parser that succeeds
func First[T any](ps ...Parser[T]) Parser[T] {
	return func(s []byte) (T, int, bool) {
		for _, p := range ps {
			res, n, ok := p(s)
			if ok {
				return res, n, true
			}
		}
		var t T
		return t, 0, false
	}
}

// Atom scans for a single atom, skipping whitespace
func Atom(val string) Parser[string] {
	b := []byte(val)
	return func(s []byte) (string, int, bool) {
		if n := ws(s); bytes.HasPrefix(s[n:], b) {
			return val, n + len(b), true
		}
		return "", 0, false
	}
}

// AtomE is like Atom but without skipping leading whitespace
func AtomE(val string) Parser[string] {
	b := []byte(val)
	return func(s []byte) (string, int, bool) {
		if bytes.HasPrefix(s, b) {
			return val, len(b), true
		}
		return "", 0, false
	}
}

// RuneIn matches a single rune of the ones in val
func RuneIn(val string) Parser[rune] {
	return func(s []byte) (rune, int, bool) {
		w := ws(s)
		res, n := utf8.DecodeRune(s[w:])
		if strings.ContainsRune(val, res) {
			return res, w + n, true
		}
		return 0, 0, false
	}
}

// Opt turns a parser into an optional parser
// Will return a result that's either Missing or the result itself
func (p Parser[T]) Opt() Parser[T] {
	return func(s []byte) (T, int, bool) {
		res, n, ok := p(s)
		if ok {
			return res, n, true
		}
		var t T
		return t, 0, true
	}
}

// Or is like Opt but returns the value given if the parser is unsuccessful
func (p Parser[T]) Or(or T) Parser[T] {
	return func(s []byte) (T, int, bool) {
		res, n, ok := p(s)
		if ok {
			return res, n, true
		}
		return or, 0, true
	}
}

// Map maps the result of a parser to a different result
// If the Mapper returns nil, the parser returns as invalid
func Map[A, B any](p Parser[A], f func(A) B) Parser[B] {
	return func(s []byte) (B, int, bool) {
		res, n, ok := p(s)
		if ok {
			return f(res), n, ok
		}
		var b B
		return b, 0, false
	}
}

// Float is a float parser
var Float = Map(Token(`[+-]?\d+(\.\d*)?([eE][+-]?\d+)?`), mapFloat)

// Int is an integer parser
var Int = Map(Token(`[+-]?\d+`), mapInt)

// Index creates a mapper that maps the result to an index in a slice
func Index[T any](i int) func([]T) T {
	return func(v []T) T {
		return v[i]
	}
}

func mapFloat(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func mapInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return i
}
