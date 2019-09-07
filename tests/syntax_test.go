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
		return left == right
	case []Token:
		for i, _ := range left.([]Token) {
			if !compareSyntax(left.([]SyntaxValue)[i], right.([]SyntaxValue)[i]) {
				return false
			}
		}
	}

	return true
}

func TestSyntax(t *testing.T) {
	var expectedSyntax, outputSyntax SyntaxValue
	var input *bufio.Reader

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
}
