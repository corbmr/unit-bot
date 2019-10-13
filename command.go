package main

import (
	"fmt"
	"strconv"

	"github.com/corbmr/unit-bot/internal/convert"
	p "github.com/corbmr/unit-bot/internal/parser"
)

type command interface {
	Do() (string, error)
}

func (c convertCommand) Do() (string, error) {
	var from convert.UnitVal
	switch uv := c.from.(type) {
	case unparsedUnitVal:
		fromUnit, err := convert.ParseUnit(string(uv.unit))
		if err != nil {
			return "", err
		}

		if fromUnit == convert.FootInch {
			fromUnit = convert.Foot
		}

		from = fromUnit.(convert.SimpleUnit).FromFloat(uv.val)

	case convert.UnitVal:
		from = uv
	}

	toUnit, err := convert.ParseUnit(string(c.to))
	if err != nil {
		return "", err
	}

	to, err := from.Convert(toUnit)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s = %s", from, to), nil
}

type unparsedUnit string

type unparsedUnitVal struct {
	val  float64
	unit unparsedUnit
}

type convertCommand struct {
	from interface{}
	to   unparsedUnit
}

var (
	intToken  = p.Token(`\d+`).Map(mapFloat)
	unitToken = p.Token(`[A-Za-z+/]+`).Map(mapUnit)
	float     = p.Token(`[+-]?\d+([.]\d*)?([eE][+-]?\d+)?`).Map(mapFloat)

	inches     = p.All(intToken, p.Atom(`"`).Opt()).Map(mapFirst)
	feetInches = p.All(intToken, p.Atom(`'`), inches.Opt()).Map(mapFeetInches)

	simpleUnitVal = p.All(float, unitToken).Map(mapSimpleUnit)

	from = p.Any(simpleUnitVal, feetInches)

	convertExpr = p.All(p.AtomE(`!conv`), from, p.Atom(`to`), unitToken).Map(mapConvertExpr)

	// Primary parser
	commandExpr = p.Any(convertExpr)
)

func mapFloat(v interface{}) interface{} {
	f, _ := strconv.ParseFloat(v.(string), 64)
	return f
}

func mapUnit(v interface{}) interface{} {
	return unparsedUnit(v.(string))
}

func mapSimpleUnit(v interface{}) interface{} {
	vs := v.([]interface{})
	return unparsedUnitVal{vs[0].(float64), vs[1].(unparsedUnit)}
}

func mapFirst(v interface{}) interface{} {
	return v.([]interface{})[0]
}

func mapFeetInches(v interface{}) interface{} {
	vs := v.([]interface{})
	feet := vs[0].(float64)
	inches := 0.0
	if f, ok := vs[2].(float64); ok {
		inches = f
	}
	return convert.FootInchVal{Feet: feet, Inches: inches}
}

func mapConvertExpr(v interface{}) interface{} {
	vs := v.([]interface{})
	return convertCommand{vs[1], vs[3].(unparsedUnit)}
}
