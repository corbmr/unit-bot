package main

import (
	"fmt"

	"github.com/corbmr/unit-bot/internal/convert"
	p "github.com/corbmr/unit-bot/internal/parser"
	"github.com/corbmr/unit-bot/internal/parser/mapper"
)

type command interface {
	Do() (string, error)
}

var errUsage = fmt.Errorf(`Usage: !conv [from][unit] to [unit]`)

func (c convertCommand) Do() (string, error) {
	if !c.ok {
		return "", errUsage
	}

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
	ok   bool
	from interface{}
	to   unparsedUnit
}

var (
	unitToken = p.Token(`[A-Za-z+/]+`).Map(mapUnit)

	inches     = p.All(p.Int, p.Atom(`"`).Opt()).Map(mapper.Index(0))
	feetInches = p.All(p.Int, p.Atom(`'`), inches.Opt()).Map(mapFeetInches)

	simpleUnitVal = p.All(p.Float, unitToken).Map(mapSimpleUnit)

	fromExpr = p.Any(simpleUnitVal, feetInches)

	// Split up like this to give error reports if the expression is incorrect
	convertExpr        = p.All(fromExpr, p.Atom(`to`), unitToken).Opt().Map(mapConvertCommand)
	convertCommandExpr = p.All(p.AtomE(`!conv`), convertExpr).Map(mapper.Index(1))

	// Primary command parser
	commandExpr = p.Any(convertCommandExpr)
)

func mapUnit(v interface{}) interface{} {
	return unparsedUnit(v.(string))
}

func mapSimpleUnit(v interface{}) interface{} {
	vs := v.([]interface{})
	return unparsedUnitVal{vs[0].(float64), vs[1].(unparsedUnit)}
}

func mapFeetInches(v interface{}) interface{} {
	vs := v.([]interface{})
	feet := vs[0].(int)
	inches := 0
	if f, ok := vs[2].(int); ok {
		inches = f
	}
	return convert.FootInchVal{Feet: float64(feet), Inches: float64(inches)}
}

func mapConvertCommand(v interface{}) interface{} {
	if c, ok := v.([]interface{}); ok {
		return convertCommand{true, c[0], c[2].(unparsedUnit)}
	}
	return convertCommand{}
}
