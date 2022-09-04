package convert

import (
	"fmt"
	"log"

	p "unit-bot/parser"
)

func Process(expr string) string {
	cmd, _, ok := convertExpr([]byte(expr))
	if !ok {
		log.Printf("Invalid command: `%v` %v\n", expr, cmd)
		return "Usage: !conv [amount][from-unit] to [to-unit]"
	}

	var from UnitVal
	switch uv := cmd.from.(type) {
	case unparsedUnitVal:
		fromUnit, ok := ParseUnit(uv.unit)
		if !ok {
			return fmt.Sprintf("Invalid unit %s", uv.unit)
		}

		from = fromUnit.FromFloat(uv.val)

	case UnitVal:
		from = uv
	}

	toUnit, ok := ParseUnit(cmd.to)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", cmd.to)
	}

	to, err := from.Convert(toUnit)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s = %s", from, to)
}

func Convert(value float64, from, to string) string {
	fromUnit, ok := ParseUnit(from)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", from)
	}

	fromValue := fromUnit.FromFloat(value)

	toUnit, ok := ParseUnit(to)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", to)
	}

	toValue, err := fromValue.Convert(toUnit)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s = %s", fromValue, toValue)
}

func Convert2(from, to string) string {
	cmd, _, ok := fromExpr([]byte(from))
	if !ok {
		log.Printf("Invalid command: `%v` %v\n", from, cmd)
		return "Usage: !conv [amount][from-unit] to [to-unit]"
	}

	var fromValue UnitVal
	switch uv := cmd.(type) {
	case unparsedUnitVal:
		fromUnit, ok := ParseUnit(uv.unit)
		if !ok {
			return fmt.Sprintf("Invalid unit %s", uv.unit)
		}

		fromValue = fromUnit.FromFloat(uv.val)

	case UnitVal:
		fromValue = uv
	}

	toUnit, ok := ParseUnit(to)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", to)
	}

	toValue, err := fromValue.Convert(toUnit)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s = %s", fromValue, toValue)
}

type unparsedUnitVal struct {
	val  float64
	unit string
}

type command struct {
	from any
	to   string
}

var (
	unitToken     = p.Token(`[A-Za-z+/$€¥£]+`)
	inches        = p.Parse2(p.Int, p.RuneIn(`"”`).Opt(), func(i int, _ rune) int { return i })
	feet          = p.Parse2(p.Int, p.RuneIn(`'’`), func(i int, _ rune) int { return i })
	feetInches    = p.Parse2(feet, inches.Or(0), mapFeetInches)
	simpleUnitVal = p.Parse2(p.Float, unitToken, mapSimpleUnit)
	currency      = p.Parse2(p.RuneIn(`$€¥£`), p.Float, mapCurrency)
	fromExpr      = p.First(simpleUnitVal, feetInches, currency)
	convertExpr   = p.Parse3(fromExpr, p.Atom(`to`), unitToken, func(v any, _ string, u string) command { return command{v, u} })
)

func mapSimpleUnit(v float64, u string) any {
	return unparsedUnitVal{v, u}
}

func mapFeetInches(feet int, inches int) any {
	return FootInchVal{Feet: float64(feet), Inches: float64(inches)}
}

func mapCurrency(c rune, v float64) any {
	return unparsedUnitVal{v, string(c)}
}
