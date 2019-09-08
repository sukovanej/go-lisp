package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

func compareSyntax(left SyntaxValue, right SyntaxValue) bool {
	switch left.(type) {
	case Token:
		switch right.(type) {
		case Token:
			return left == right
		default:
			return false
		}
	case []SyntaxValue:
		switch right.(type) {
		case []SyntaxValue:
			if len(left.([]SyntaxValue)) != len(right.([]SyntaxValue)) {
				return false
			}

			for i, _ := range left.([]SyntaxValue) {
				if !compareSyntax(left.([]SyntaxValue)[i], right.([]SyntaxValue)[i]) {
					return false
				} else {
					return true
				}
			}
		default:
			return false
		}
	case nil:
		return left == right
	}

	return false
}

var expectedSyntax, outputSyntax SyntaxValue
var input *bufio.Reader

func TestSyntaxBasic(t *testing.T) {
	input = bufio.NewReader(strings.NewReader("(test)"))
	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"test", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	input = bufio.NewReader(strings.NewReader("(test value)"))
	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"test", SYMBOL}),
			SyntaxValue(Token{"value", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxWithInner(t *testing.T) {
	input = bufio.NewReader(strings.NewReader("(test (fn 1 2) value)"))
	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"test", SYMBOL}),
			SyntaxValue(
				[]SyntaxValue{
					SyntaxValue(Token{"fn", SYMBOL}),
					SyntaxValue(Token{"1", SYMBOL}),
					SyntaxValue(Token{"2", SYMBOL}),
				},
			),
			SyntaxValue(Token{"value", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	input = bufio.NewReader(strings.NewReader("(+ 1 2)"))
	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"+", SYMBOL}),
			SyntaxValue(Token{"1", SYMBOL}),
			SyntaxValue(Token{"2", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	input = bufio.NewReader(strings.NewReader("(+ (+ 1 2) 3)"))
	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"+", SYMBOL}),
			SyntaxValue(
				[]SyntaxValue{
					SyntaxValue(Token{"+", SYMBOL}),
					SyntaxValue(Token{"1", SYMBOL}),
					SyntaxValue(Token{"2", SYMBOL}),
				},
			),
			SyntaxValue(Token{"3", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	input = bufio.NewReader(strings.NewReader("((fn (x) (+ x 1)) 1)"))
	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue([]SyntaxValue{
				SyntaxValue(Token{"fn", SYMBOL}),
				SyntaxValue(
					[]SyntaxValue{
						SyntaxValue(Token{"x", SYMBOL}),
					},
				),
				SyntaxValue(
					[]SyntaxValue{
						SyntaxValue(Token{"+", SYMBOL}),
						SyntaxValue(Token{"x", SYMBOL}),
						SyntaxValue(Token{"1", SYMBOL}),
					},
				),
			}),
			SyntaxValue(Token{"1", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxMultipleStatements(t *testing.T) {
	input = bufio.NewReader(strings.NewReader("(test (fn 1 2) value)\n(var x)"))

	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"test", SYMBOL}),
			SyntaxValue(
				[]SyntaxValue{
					SyntaxValue(Token{"fn", SYMBOL}),
					SyntaxValue(Token{"1", SYMBOL}),
					SyntaxValue(Token{"2", SYMBOL}),
				},
			),
			SyntaxValue(Token{"value", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	expectedSyntax = SyntaxValue(
		[]SyntaxValue{
			SyntaxValue(Token{"var", SYMBOL}),
			SyntaxValue(Token{"x", SYMBOL}),
		},
	)
	outputSyntax = GetSyntax(input)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	expectedSyntax = nil
	outputSyntax = GetSyntax(input)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}
