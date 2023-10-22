package convert

import (
	"fmt"
	"log/slog"
	"strings"

	p "unit-bot/parser"
)

func Process(expr string) string {
	cmd, _, ok := convertExpr([]byte(expr))
	if !ok {
		slog.Info("Invalid command: `%v` %v\n", expr, cmd)
		return "Usage: !conv [amount][from-unit] to [to-unit]"
	}

	var from UnitVal
	switch uv := cmd.from.(type) {
	case unparsedUnitVal:
		fromUnit, ok := LookupUnit(uv.unit)
		if !ok {
			return fmt.Sprintf("Invalid unit %s", uv.unit)
		}

		from = fromUnit.FromFloat(uv.val)

	case UnitVal:
		from = uv
	}

	toUnit, ok := LookupUnit(cmd.to)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", cmd.to)
	}

	slog.Debug("converting", "from", from, "to", toUnit)

	to, err := from.Convert(toUnit)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s = %s", from, to)
}

func Convert(from, to string) string {
	cmd, _, ok := fromExpr([]byte(from))
	if !ok {
		slog.Info("Invalid command: `%v` %v\n", from, cmd)
		return "Usage: !conv [amount][from-unit] to [to-unit]"
	}

	var fromValue UnitVal
	switch uv := cmd.(type) {
	case unparsedUnitVal:
		fromUnit, ok := LookupUnit(uv.unit)
		if !ok {
			return fmt.Sprintf("Invalid unit %s", uv.unit)
		}

		fromValue = fromUnit.FromFloat(uv.val)

	case UnitVal:
		fromValue = uv
	}

	toUnit, ok := LookupUnit(to)
	if !ok {
		return fmt.Sprintf("Invalid unit %s", to)
	}

	slog.Debug("converting", "from", debug(fromValue), "to", debug(toUnit))

	toValue, err := fromValue.Convert(toUnit)
	if err != nil {
		slog.Error("Cannot convert",
			"fromValue", debug(fromValue), "toUnit", debug(toUnit), "err", err)
		return err.Error()
	}

	return fmt.Sprintf("%s = %s", fromValue, toValue)
}

func debug(v any) string {
	return fmt.Sprintf("%#v", v)
}

func Autocomplete(from, to string) []string {
	cmd, _, ok := fromExpr([]byte(from))
	if !ok {
		slog.Info("Invalid command: `%v` %v\n", from, cmd)
		return nil
	}

	var fromValue UnitVal
	switch uv := cmd.(type) {
	case unparsedUnitVal:
		fromUnit, ok := LookupUnit(uv.unit)
		if !ok {
			return nil
		}

		fromValue = fromUnit.FromFloat(uv.val)

	case UnitVal:
		fromValue = uv
	}

	options := []string{}
	for _, unit := range unitDimensionMap[fromValue.Unit().Dimension()] {
		if strings.HasPrefix(unit.String(), to) {
			options = append(options, unit.String())
		}
	}

	return options
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
	inches        = p.Parse2(p.Int, p.RuneIn(`"”`).Opt(), fst[int, rune])
	feet          = p.Parse2(p.Int, p.RuneIn(`'’`), fst[int, rune])
	feetInches    = p.Parse2(feet, inches.Or(0), mapFeetInches)
	simpleUnitVal = p.Parse2(p.Float, unitToken, mapSimpleUnit)
	currency      = p.Parse2(p.RuneIn(`$€¥£`), p.Float, mapCurrency)
	fromExpr      = p.First(simpleUnitVal, feetInches, currency)
	convertExpr   = p.Parse3(fromExpr, p.Atom(`to`), unitToken, func(v any, _ string, u string) command { return command{v, u} })
)

func fst[A any, B any](a A, b B) A {
	return a
}

func mapSimpleUnit(v float64, u string) any {
	return unparsedUnitVal{v, u}
}

func mapFeetInches(feet int, inches int) any {
	return FootInchVal{Feet: float64(feet), Inches: float64(inches)}
}

func mapCurrency(c rune, v float64) any {
	return unparsedUnitVal{v, string(c)}
}
