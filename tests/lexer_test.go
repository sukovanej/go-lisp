package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

var expectedToken Token
var inputLexer *bufio.Reader

func TestToken(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader("(test)"))
	expectedResult := []Token{
		Token{"(", TOKEN_LPAR},
		Token{"test", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	inputLexer = bufio.NewReader(strings.NewReader("(symbol-1 symbol-2 (hello world))"))
	expectedResult = []Token{
		Token{"(", TOKEN_LPAR},
		Token{"symbol-1", SYMBOL},
		Token{"symbol-2", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"hello", SYMBOL},
		Token{"world", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	inputLexer = bufio.NewReader(
		strings.NewReader("(symbol-1 symbol-2 (hello world (fn 1 2)) another terms)"),
	)
	expectedResult = []Token{
		Token{"(", TOKEN_LPAR},
		Token{"symbol-1", SYMBOL},
		Token{"symbol-2", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"hello", SYMBOL},
		Token{"world", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"fn", SYMBOL},
		Token{"1", SYMBOL},
		Token{"2", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{"another", SYMBOL},
		Token{"terms", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenMultipleStatements(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader(" (test  (fn   1 2) value)  \n (var   x)  "))

	expectedResult := []Token{
		Token{"(", TOKEN_LPAR},
		Token{"test", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"fn", SYMBOL},
		Token{"1", SYMBOL},
		Token{"2", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"value", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"(", TOKEN_LPAR},
		Token{"var", SYMBOL},
		Token{"x", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	inputLexer = bufio.NewReader(strings.NewReader(`(set fact (fn (x) (
		if (== x 1) 
			1 
			(* x (fact (- x 1))))))
		(fact 4)`))
	expectedResult = []Token{
		Token{"(", TOKEN_LPAR},
		Token{"set", SYMBOL},
		Token{"fact", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"fn", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"x", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"(", TOKEN_LPAR},
		Token{"if", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"==", SYMBOL},
		Token{"x", SYMBOL},
		Token{"1", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"1", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"*", SYMBOL},
		Token{"x", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"fact", SYMBOL},
		Token{"(", TOKEN_LPAR},
		Token{"-", SYMBOL},
		Token{"x", SYMBOL},
		Token{"1", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{")", TOKEN_RPAR},
		Token{"(", TOKEN_LPAR},
		Token{"fact", SYMBOL},
		Token{"4", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenString(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader("(test \"asdf asdf\" 1)"))
	expectedResult := []Token{
		Token{"(", TOKEN_LPAR},
		Token{"test", SYMBOL},
		Token{"asdf asdf", SYMBOL_STRING},
		Token{"1", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}
