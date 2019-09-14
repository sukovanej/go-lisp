package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

func compareSyntax(left SyntaxValue, right SyntaxValue) bool {
	if left == nil || right == nil {
		return left == right
	}

	switch left.GetType() {
	case SYNTAX_SYMBOL:
		switch right.GetType() {
		case SYNTAX_SYMBOL:

			return left == right
		default:
			return false
		}
	case SYNTAX_LIST:
		switch right.GetType() {
		case SYNTAX_LIST:
			leftValue := left.(ListValue).Value
			rightValue := right.(ListValue).Value

			if len(leftValue) != len(rightValue) {
				return false
			}

			for i := 0; i < len(leftValue); i++ {
				if !compareSyntax(leftValue[i], rightValue[i]) {
					return false
				}
			}

			return true
		default:
			return false
		}
	case SYNTAX_LIST_LITERAL:
		switch right.GetType() {
		case SYNTAX_LIST_LITERAL:
			leftValue := left.(ListLiteralValue).Value
			rightValue := right.(ListLiteralValue).Value

			if len(leftValue) != len(rightValue) {
				return false
			}

			for i := 0; i < len(leftValue); i++ {
				if !compareSyntax(leftValue[i], rightValue[i]) {
					return false
				}
			}
			return true

		default:
			return false
		}
	case SYNTAX_DICT_LITERAL:
		switch right.GetType() {
		case SYNTAX_DICT_LITERAL:
			leftValue := left.(DictLiteralValue).Value
			rightValue := right.(DictLiteralValue).Value

			if len(leftValue) != len(rightValue) {
				return false
			}

			for i := 0; i < len(leftValue); i++ {
				if !compareSyntax(leftValue[i], rightValue[i]) {
					return false
				}
			}
			return true

		default:
			return false
		}
	}

	return false
}

var expectedSyntax, outputSyntax SyntaxValue
var inputSyntax *bufio.Reader

func TestSyntaxBasic(t *testing.T) {
	inputSyntax = bufio.NewReader(strings.NewReader("(test)"))
	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	inputSyntax = bufio.NewReader(strings.NewReader("(test value)"))
	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
			SymbolValue{Token{"value", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxWithInner(t *testing.T) {
	inputSyntax = bufio.NewReader(strings.NewReader("(test (fn 1 2) value)"))
	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
			ListValue{
				[]SyntaxValue{
					SymbolValue{Token{"fn", SYMBOL}},
					SymbolValue{Token{"1", SYMBOL}},
					SymbolValue{Token{"2", SYMBOL}},
				},
			},
			SymbolValue{Token{"value", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	inputSyntax = bufio.NewReader(strings.NewReader("(+ 1 2)"))
	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"+", SYMBOL}},
			SymbolValue{Token{"1", SYMBOL}},
			SymbolValue{Token{"2", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)

	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	inputSyntax = bufio.NewReader(strings.NewReader("(+ (+ 1 2) 3)"))
	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"+", SYMBOL}},
			ListValue{
				[]SyntaxValue{
					SymbolValue{Token{"+", SYMBOL}},
					SymbolValue{Token{"1", SYMBOL}},
					SymbolValue{Token{"2", SYMBOL}},
				},
			},
			SymbolValue{Token{"3", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	inputSyntax = bufio.NewReader(strings.NewReader("((fn (x) (+ x 1)) 1)"))
	expectedSyntax = ListValue{
		[]SyntaxValue{
			ListValue{
				[]SyntaxValue{
					SymbolValue{Token{"fn", SYMBOL}},
					ListValue{
						[]SyntaxValue{
							SymbolValue{Token{"x", SYMBOL}},
						},
					},
					ListValue{
						[]SyntaxValue{
							SymbolValue{Token{"+", SYMBOL}},
							SymbolValue{Token{"x", SYMBOL}},
							SymbolValue{Token{"1", SYMBOL}},
						},
					},
				},
			},
			SymbolValue{Token{"1", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxMultipleStatements(t *testing.T) {
	inputSyntax = bufio.NewReader(strings.NewReader("(test (fn 1 2) value)\n(var x)"))

	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
			ListValue{
				[]SyntaxValue{
					SymbolValue{Token{"fn", SYMBOL}},
					SymbolValue{Token{"1", SYMBOL}},
					SymbolValue{Token{"2", SYMBOL}},
				},
			},
			SymbolValue{Token{"value", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"var", SYMBOL}},
			SymbolValue{Token{"x", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	expectedSyntax = nil
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxString(t *testing.T) {
	inputSyntax = bufio.NewReader(strings.NewReader(`(test "hello world" 1)`))

	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
			SymbolValue{Token{"hello world", SYMBOL_STRING}},
			SymbolValue{Token{"1", SYMBOL}},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxList(t *testing.T) {
	inputSyntax = bufio.NewReader(strings.NewReader(`(test [1 2 3 4])`))

	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
			ListLiteralValue{
				[]SyntaxValue{
					SymbolValue{Token{"1", SYMBOL}},
					SymbolValue{Token{"2", SYMBOL}},
					SymbolValue{Token{"3", SYMBOL}},
					SymbolValue{Token{"4", SYMBOL}},
				},
			},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}

	inputSyntax = bufio.NewReader(strings.NewReader(`[]`))

	expectedSyntax = ListLiteralValue{
		[]SyntaxValue{},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

func TestSyntaxDict(t *testing.T) {
	inputSyntax = bufio.NewReader(strings.NewReader(`(test {"key" "value"})`))

	expectedSyntax = ListValue{
		[]SyntaxValue{
			SymbolValue{Token{"test", SYMBOL}},
			DictLiteralValue{
				[]SyntaxValue{
					SymbolValue{Token{"key", SYMBOL_STRING}},
					SymbolValue{Token{"value", SYMBOL_STRING}},
				},
			},
		},
	}
	outputSyntax = GetSyntax(inputSyntax)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}
