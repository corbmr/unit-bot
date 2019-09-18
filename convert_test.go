package main

import "testing"

func TestRegex(t *testing.T) {

	should := func(s string) {
		m := inputRegex.FindString(s)
		if m == "" {
			t.Errorf("Regex expected to match: {%s}", s)
		}
	}

	shouldnt := func(s string) {
		m := inputRegex.FindString(s)
		if m != "" {
			t.Errorf("Regex not expected to match: {%s}. Found {%s}", s, m)
		}
	}

	should(`6.4cm to m`)
	should(`6' to m`)
	should(`6'0 to m`)
	should(`5 ft to m`)
	shouldnt(`5. to m`)
	should(`5. ft to m`)
	shouldnt(``)
}
