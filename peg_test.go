package main

import (
	"testing"
)

func genericTest[T comparable](t *testing.T, parser Parser[T], source string, expected ParseResult[T], loc int) {
	result, loc2 := parser([]byte(source), 0)
	if result != expected {
		t.Errorf("parser(%s, 0) = %v; want %v", source, result, expected)
	}
	if loc2 != loc {
		t.Errorf("loc = %d; want %d", loc2, loc)
	}
}

func TestChar(t *testing.T) {
	genericTest(t, char('a'), "abc", succeed(byte('a')), 1)
}

func TestDigit(t *testing.T) {
	genericTest(t, digit(), "123", succeed(1), 1)
}

func TestSeqAndApply(t *testing.T) {
	source := "abcd"
	f := func(bs ParseResult[[]byte]) ParseResult[string] {
		var res ParseResult[string]
		res.success = false
		if bs.success {
			res.val = string(bs.val)
			res.success = true
		}
		return res
	}
	parser := seq[byte]([]Parser[byte]{char('a'), char('b'), char('c')})
	result := succeed("abc")
	genericTest(t, apply(f, parser), source, result, 3)
}

func TestChoice(t *testing.T) {
	parser := choice[byte]([]Parser[byte]{char('a'), char('b')})
	genericTest(t, parser, "abc", succeed(byte('a')), 1)
	genericTest(t, parser, "bbc", succeed(byte('b')), 1)
}

func TestMany(t *testing.T) {
	source := "123abc"
	parser := many(digit())
	expected := succeed([]int{1, 2, 3})
	result, loc := parser([]byte(source), 0)
	if !result.success {
		t.Errorf("result.success = false; expected true")
	}
	if len(result.val) != len(expected.val) {
		t.Errorf("len(res) = %d; want %d", len(result.val), len(expected.val))
	}
	for i := 0; i < len(expected.val); i++ {
		if result.val[i] != expected.val[i] {
			t.Errorf("result = %v; want %v", result, expected)
		}
	}
	if loc != 3 {
		t.Errorf("loc = %d; want 3", loc)
	}
}

func TestMany1(t *testing.T) {
	parser := many1(char('a'))
	result, loc := parser([]byte("bcd"), 0)
	if result.success {
		t.Errorf("result.success = true; want false")
	}

	result, loc = parser([]byte("aaabcd"), 0)
	expected := succeed([]byte{'a', 'a', 'a'})
	if !result.success {
		t.Errorf("result.success = false; want true")
	}
	if len(result.val) != len(expected.val) {
		t.Errorf("result = %v; want %v", result, expected)
	}
	for i := 0; i < len(expected.val); i++ {
		if result.val[i] != expected.val[i] {
			t.Errorf("result = %v; want %v", result, expected)
		}
	}
	if loc != 3 {
		t.Errorf("loc = %d; want 3", loc)
	}
}
