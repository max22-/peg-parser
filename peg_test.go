package main

import (
	"testing"
)

func genericTest(t *testing.T, parser Parser, source string, expected ParseResult, loc int) {
	result, loc2 := parser([]byte(source), 0)
	if result != expected {
		t.Errorf("parser(%s, 0) = %v; want %v", source, result, expected)
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

func TestSeqAndApply(t *testing.T) {
	source := "abcd"
	f := func(r ParseResult) ParseResult {
		res := ""
		for _, c := range r.(List) {
			res += string(c.(Char))
		}
		return String(res)
	}
	parser := seq([]Parser{char('a'), char('b'), char('c')})
	result := String("abc")
	genericTest(t, apply(f, parser), source, result, 3)
}

func TestChoice(t *testing.T) {
	parser := choice([]Parser{char('a'), digit()})
	genericTest(t, parser, "abc", Char('a'), 1)
	genericTest(t, parser, "1bc", Int(1), 1)
}

func TestMany(t *testing.T) {
	source := "123abc"
	parser := many(digit())
	expected := List{Int(1), Int(2), Int(3)}
	result, loc := parser([]byte(source), 0)
	if len(result.(List)) != len(expected) {
		t.Errorf("len(res) = %d; want %d", len(result.(List)), len(expected))
	}
	for i := 0; i < len(expected); i++ {
		if result.(List)[i] != expected[i] {
			t.Errorf("result = %v; want %v", result, expected)
		}
	}
	if loc != 3 {
		t.Errorf("loc = %d; want 3", loc)
	}
}

func TestMany1(t *testing.T) {
	parser := many1(char('a'))
	genericTest(t, parser, "bcd", nil, 0)

	result, loc := parser([]byte("aaabcd"), 0)
	expected := List{Char('a'), Char('a'), Char('a')}
	if len(result.(List)) != len(expected) {
		t.Errorf("len(result) = %d; want %d", len(result.(List)), len(expected))
	}
	for i := 0; i < len(expected); i++ {
		if result.(List)[i] != expected[i] {
			t.Errorf("result = %v; want %v", result, expected)
		}
	}
	if loc != 3 {
		t.Errorf("loc = %d; want 3", loc)
	}
}
