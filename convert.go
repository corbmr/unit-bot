package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const usage = `Usage: !conv [from][unit] to [unit]`

// TODO: write a real parser, this is gross
var inputRegex = regexp.MustCompile(
	`^(((?P<num>[+-]?\d+([.]\d*)?([eE][+-]?\d+)?)\s*(?P<from>[A-Za-z+]+))|((?P<feet>[+-]?\d+)'\s*((?P<inches>\d+)"?)?))\s+to\s+(?P<to>[A-Za-z+]+)`)

func generateResponse(inp string) (string, error) {
	match := findNamed(inputRegex, inp)
	if match == nil {
		return "", fmt.Errorf(usage)
	}

	toUnit, err := parseUnit(match["to"])
	if err != nil {
		return "", err
	}

	var from unitVal

	switch {
	case match["num"] != "" && match["from"] != "":
		log.Printf("num: [%s], from: [%s]", match["num"], match["from"])
		fromUnit, err := parseUnit(match["from"])
		if err != nil {
			return "", err
		}

		if fromUnit == footInch {
			fromUnit = foot
		}

		// This shouldn't fail since the regex guaranteed we got a simple unit
		simple := fromUnit.(simpleUnit)

		// Again, shouldn't fail because of regex
		float, _ := strconv.ParseFloat(match["num"], 64)

		from = simple.fromFloat(float)

	case match["feet"] != "":
		feet, _ := strconv.ParseFloat(match["feet"], 64)

		var inches float64
		if match["inches"] != "" {
			inches, _ = strconv.ParseFloat(match["inches"], 64)
		}

		from = footInchVal{feet, inches}
	}

	to, err := from.convert(toUnit)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s = %s", from, to), nil
}

func findNamed(r *regexp.Regexp, s string) map[string]string {
	match := r.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	names := r.SubexpNames()
	out := make(map[string]string)
	for i, n := range names {
		if n != "" {
			out[n] = match[i]
		}
	}

	return out
}

func parseUnit(s string) (unitType, error) {
	u, ok := unitMap[strings.ToLower(s)]
	if !ok {
		return nil, fmt.Errorf("Invalid unit %s", s)
	}
	return u, nil
}
