package main

import (
	"fmt"

	"unit-bot/internal/convert"
	p "unit-bot/parser"
)

type command interface {
	Do() (string, error)
}

// ErrorUsage occurs when formatting is incorrect
var ErrorUsage = fmt.Errorf(`Usage: !conv [from][unit] to [unit]`)

// ErrorInvalidUnit occurs when an invalid unit is given
type ErrorInvalidUnit string

func (err ErrorInvalidUnit) Error() string {
	return fmt.Sprintf("Invalid unit %s", string(err))
}

func (c convertCommand) Do() (string, error) {
	if !c.ok {
		return "", ErrorUsage
	}

	var from convert.UnitVal
	switch uv := c.from.(type) {
	case unparsedUnitVal:
		fromUnit, ok := convert.ParseUnit(uv.unit)
		if !ok {
			return "", ErrorInvalidUnit(uv.unit)
		}

		from = fromUnit.FromFloat(uv.val)

	case convert.UnitVal:
		from = uv
	}

	toUnit, ok := convert.ParseUnit(c.to)
	if !ok {
		return "", ErrorInvalidUnit(c.to)
	}

	to, err := from.Convert(toUnit)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s = %s", from, to), nil
}

type unparsedUnitVal struct {
	val  float64
	unit string
}

type convertCommand struct {
	ok   bool
	from interface{}
	to   string
}

var (
	unitToken = p.Token(`[A-Za-z+/$€¥£]+`)

	inches     = p.All(p.Int, p.Atom(`"`).Opt()).Map(p.Index(0))
	feetInches = p.All(p.Int, p.Atom(`'`), inches.Or(0)).Map(mapFeetInches)

	simpleUnitVal = p.All(p.Float, unitToken).Map(mapSimpleUnit)

	currency = p.All(p.RuneIn(`$€¥£`), p.Float).Map(mapCurrency)

	fromExpr = p.Any(simpleUnitVal, feetInches, currency)

	convertPrefix = p.Token(`(?i)!conv(ert)?`)

	// Split up like this to give error reports if the expression is incorrect
	convertExpr        = p.All(fromExpr, p.Atom(`to`), unitToken).Map(mapConvertCommand)
	convertCommandExpr = p.All(convertPrefix, convertExpr.Or(convertCommand{})).Map(p.Index(1))

	// Primary command parser
	commandExpr = p.Any(convertCommandExpr)
)

func mapSimpleUnit(v interface{}) interface{} {
	vs := v.([]interface{})
	return unparsedUnitVal{vs[0].(float64), vs[1].(string)}
}

func mapFeetInches(v interface{}) interface{} {
	vs := v.([]interface{})
	feet := vs[0].(int)
	inches := vs[2].(int)
	return convert.FootInchVal{Feet: float64(feet), Inches: float64(inches)}
}

func mapConvertCommand(v interface{}) interface{} {
	c := v.([]interface{})
	return convertCommand{true, c[0], c[2].(string)}
}

func mapCurrency(v interface{}) interface{} {
	c := v.([]interface{})
	u, ok := convert.ParseUnit(string(c[0].(rune)))
	if !ok {
		return nil
	}
	return convert.CurrencyVal{V: c[1].(float64), U: u.(*convert.CurrencyUnit)}
}
