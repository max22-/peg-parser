package pegparser

import (
	"reflect"
	"testing"
)

func assertEqual[T comparable](t *testing.T, name string, val T, expected T) {
	if val != expected {
		t.Errorf("%s = %v; want %v", name, val, expected)
	}
}

func assertEqualSlice[T comparable](t *testing.T, name string, val []T, expected []T) {
	if !reflect.DeepEqual(val, expected) {
		t.Errorf("%s = %v; want %v", name, val, expected)
	}
}

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
	genericTest(t, Char('a'), "abc", succeed(byte('a')), 1)
}

func TestDigit(t *testing.T) {
	genericTest(t, Digit(), "123", succeed(1), 1)
}

func TestSeqAndApply(t *testing.T) {
	source := "abcd"
	parser := Seq[byte]([]Parser[byte]{Char('a'), Char('b'), Char('c')})
	result := succeed("abc")
	genericTest(t, Apply(func(bs []byte) string { return string(bs) }, parser), source, result, 3)
}

func TestChoice(t *testing.T) {
	parser := Choice[byte]([]Parser[byte]{Char('a'), Char('b')})
	genericTest(t, parser, "abc", succeed(byte('a')), 1)
	genericTest(t, parser, "bbc", succeed(byte('b')), 1)
}

func TestMany(t *testing.T) {
	source := "123abc"
	parser := Many(Digit())
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
	parser := Many1(Char('a'))
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

func TestMaybe(t *testing.T) {
	parser := Maybe(Char('a'))
	result, loc := parser([]byte("bcd"), 0)
	if !result.success {
		t.Error("result.success = false; want true")
	}
	if result.val != nil {
		t.Errorf("result.val = %v; want nil", result.val)
	}
	if loc != 0 {
		t.Errorf("loc = %d; want 0", loc)
	}

	result, loc = parser([]byte("abcd"), 0)
	expected := byte('a')
	if !result.success {
		t.Error("result.success = false; want true")
	}
	if *result.val != expected {
		t.Errorf("result.val = %v; want %v", result.val, expected)
	}
	if loc != 1 {
		t.Errorf("loc = %d; want 1", loc)
	}
}

func TestAnd(t *testing.T) {
	parser := And(Char('a'), Char('b'))
	result, loc := parser([]byte("abc"), 0)
	assertEqual(t, "result.success", result.success, true)
	assertEqual(t, "result.val", result.val, 'a')
	assertEqual(t, "loc", loc, 1)

	result, loc = parser([]byte("acb"), 0)
	var zeroChar byte
	assertEqual(t, "result.success", result.success, false)
	assertEqual(t, "result.val", result.val, zeroChar)
	assertEqual(t, "loc", loc, 0)

	result, loc = parser([]byte("bbc"), 0)
	assertEqual(t, "result.success", result.success, false)
	assertEqual(t, "result.val", result.val, zeroChar)
	assertEqual(t, "loc", loc, 0)
}

func TestNot(t *testing.T) {
	parser := Not(Char('a'))
	var zeroChar byte
	result, loc := parser([]byte("abc"), 0)
	assertEqual(t, "result.success", result.success, false)
	assertEqual(t, "result.val", result.val, zeroChar)
	assertEqual(t, "loc", loc, 0)

	result, loc = parser([]byte("bbc"), 0)
	assertEqual(t, "result.success", result.success, true)
	assertEqual(t, "result.val", result.val, zeroChar)
	assertEqual(t, "loc", loc, 0)
}
