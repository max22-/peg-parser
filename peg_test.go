package main

import "testing"

func genericTest(t *testing.T, parser Parser, source string, result ParseResult, loc int) {
	res, loc2 := parser([]byte(source), 0)
	if res != result {
		t.Errorf("parser(%s, 0) = %v; want %v", source, res, result)
	}
	if loc2 != loc {
		t.Errorf("loc = %d; want %d", loc2, loc)
	}
}

func TestChar(t *testing.T) {
	genericTest(t, char('a'), "abc", Char('a'), 1)
}

func TestDigit(t *testing.T) {
	genericTest(t, digit(), "123", Int(1), 1)
}
