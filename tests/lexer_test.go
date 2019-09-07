package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

func TestToken(t *testing.T) {
	var expectedToken Token
	var input *bufio.Reader

	input = bufio.NewReader(strings.NewReader("(test)"))
	expectedResult := []Token{
		Token{"(", TOKEN_LPAR},
		Token{"test", SYMBOL},
		Token{")", TOKEN_RPAR},
		Token{"", END},
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(input)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	input = bufio.NewReader(strings.NewReader("(symbol-1 symbol-2 (hello world))"))
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
		expectedToken = GetToken(input)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	input = bufio.NewReader(
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
		expectedToken = GetToken(input)
		if token != expectedToken {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}
