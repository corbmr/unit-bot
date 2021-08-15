package convert

import (
	"fmt"
	"log"

	p "unit-bot/parser"
)

func Process(cmd string) string {
	res, _, ok := convertExpr([]byte(cmd))
	if !ok {
		log.Printf("Invalid command: `%v` %v\n", cmd, res)
		return "Usage: !conv [amount][from-unit] to [to-unit]"
	}

	c := res.([]interface{})
	cmdFrom := c[0]
	cmdTo := c[2].(string)

	var from UnitVal
	switch uv := cmdFrom.(type) {
	case unparsedUnitVal:
		fromUnit, ok := ParseUnit(uv.unit)
		if !ok {
			return fmt.Sprintf("Invalid unit %s", uv.unit)
		}

		from = fromUnit.FromFloat(uv.val)

	case UnitVal:
		from = uv
	}

	toUnit, ok := ParseUnit(cmdTo)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", cmdTo)
	}

	to, err := from.Convert(toUnit)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s = %s", from, to)
}

type unparsedUnitVal struct {
	val  float64
	unit string
}

var (
	unitToken = p.Token(`[A-Za-z+/$€¥£]+`)

	inches     = p.All(p.Int, p.Atom(`"`).Opt()).Map(p.Index(0))
	feetInches = p.All(p.Int, p.Atom(`'`), inches.Or(0)).Map(mapFeetInches)

	simpleUnitVal = p.All(p.Float, unitToken).Map(mapSimpleUnit)

	currency = p.All(p.RuneIn(`$€¥£`), p.Float).Map(mapCurrency)

	fromExpr = p.Any(simpleUnitVal, feetInches, currency)

	convertExpr = p.All(fromExpr, p.Atom(`to`), unitToken)
)

func mapSimpleUnit(v interface{}) interface{} {
	vs := v.([]interface{})
	return unparsedUnitVal{vs[0].(float64), vs[1].(string)}
}

func mapFeetInches(v interface{}) interface{} {
	vs := v.([]interface{})
	feet := vs[0].(int)
	inches := vs[2].(int)
	return FootInchVal{Feet: float64(feet), Inches: float64(inches)}
}

func mapCurrency(v interface{}) interface{} {
	c := v.([]interface{})
	u, ok := ParseUnit(string(c[0].(rune)))
	if !ok {
		return nil
	}
	return CurrencyVal{V: c[1].(float64), U: u.(*CurrencyUnit)}
}
