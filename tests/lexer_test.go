package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

var expectedToken Token
var inputLexer *bufio.Reader

func CreateToken(s string, t TokenType) Token {
	return Token{s, t, 0, 0, ""}
}

func CompareTokens(t1 Token, t2 Token) bool {
	return t1.Type == t2.Type && t1.Symbol == t2.Symbol
}

func TestToken(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader("(test)"))
	expectedResult := []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("test", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	inputLexer = bufio.NewReader(strings.NewReader("(symbol-1 symbol-2 (hello world))"))
	expectedResult = []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("symbol-1", SYMBOL),
		CreateToken("symbol-2", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("hello", SYMBOL),
		CreateToken("world", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	inputLexer = bufio.NewReader(
		strings.NewReader("(symbol-1 symbol-2 (hello world (fn 1 2)) another terms)"),
	)
	expectedResult = []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("symbol-1", SYMBOL),
		CreateToken("symbol-2", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("hello", SYMBOL),
		CreateToken("world", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("fn", SYMBOL),
		CreateToken("1", SYMBOL),
		CreateToken("2", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("another", SYMBOL),
		CreateToken("terms", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenMultipleStatements(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader(" (test  (fn   1 2) value)  \n (var   x)  "))

	expectedResult := []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("test", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("fn", SYMBOL),
		CreateToken("1", SYMBOL),
		CreateToken("2", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("value", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("var", SYMBOL),
		CreateToken("x", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}

	inputLexer = bufio.NewReader(strings.NewReader(`(set fact (fn (x) (
		if (== x 1) 
			1 
			(* x (fact (- x 1))))))
		(fact 4)`))
	expectedResult = []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("set", SYMBOL),
		CreateToken("fact", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("fn", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("x", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("if", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("==", SYMBOL),
		CreateToken("x", SYMBOL),
		CreateToken("1", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("1", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("*", SYMBOL),
		CreateToken("x", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("fact", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("-", SYMBOL),
		CreateToken("x", SYMBOL),
		CreateToken("1", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("fact", SYMBOL),
		CreateToken("4", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenString(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader("(test \"asdf asdf\" 1)"))
	expectedResult := []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("test", SYMBOL),
		CreateToken("asdf asdf", SYMBOL_STRING),
		CreateToken("1", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenStruct(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader(`
(set person (struct
    name
    age
))

(set-> person name "Adam")
(print (-> person name))
`))
	expectedResult := []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("set", SYMBOL),
		CreateToken("person", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("struct", SYMBOL),
		CreateToken("name", SYMBOL),
		CreateToken("age", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("set->", SYMBOL),
		CreateToken("person", SYMBOL),
		CreateToken("name", SYMBOL),
		CreateToken("Adam", SYMBOL_STRING),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("print", SYMBOL),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("->", SYMBOL),
		CreateToken("person", SYMBOL),
		CreateToken("name", SYMBOL),
		CreateToken(")", TOKEN_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenDictAndList(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader(`
(set l [1 2 3])
(set d {"key" "value"})
`))
	expectedResult := []Token{
		CreateToken("(", TOKEN_LPAR),
		CreateToken("set", SYMBOL),
		CreateToken("l", SYMBOL),
		CreateToken("[", TOKEN_LIST_LPAR),
		CreateToken("1", SYMBOL),
		CreateToken("2", SYMBOL),
		CreateToken("3", SYMBOL),
		CreateToken("]", TOKEN_LIST_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("(", TOKEN_LPAR),
		CreateToken("set", SYMBOL),
		CreateToken("d", SYMBOL),
		CreateToken("{", TOKEN_DICT_LPAR),
		CreateToken("key", SYMBOL_STRING),
		CreateToken("value", SYMBOL_STRING),
		CreateToken("}", TOKEN_DICT_RPAR),
		CreateToken(")", TOKEN_RPAR),
		CreateToken("", END),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}

func TestTokenSpecialStringCharacters(t *testing.T) {
	inputLexer = bufio.NewReader(strings.NewReader(`"\t\n\""`))
	expectedResult := []Token{
		CreateToken("\t\n\"", SYMBOL_STRING),
	}

	for _, token := range expectedResult {
		expectedToken = GetToken(inputLexer, &BufferMetaInformation{0, 0, "<test>"})
		if !CompareTokens(token, expectedToken) {
			t.Errorf("%v != %v.", token, expectedToken)
		}
	}
}
