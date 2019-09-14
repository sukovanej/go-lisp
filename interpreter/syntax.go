package interpreter

import (
	"bufio"
)

type SyntaxType int

const (
	SYNTAX_SYMBOL SyntaxType = iota
	SYNTAX_LIST
	SYNTAX_LIST_LITERAL
	SYNTAX_DICT_LITERAL
)

type SyntaxValue interface {
	GetType() SyntaxType
}

type SymbolValue struct {
	Value Token
}

func (s SymbolValue) GetType() SyntaxType {
	return SYNTAX_SYMBOL
}

type ListValue struct {
	Value []SyntaxValue
}

func (_ ListValue) GetType() SyntaxType {
	return SYNTAX_LIST
}

type ListLiteralValue struct {
	Value []SyntaxValue
}

func (_ ListLiteralValue) GetType() SyntaxType {
	return SYNTAX_LIST_LITERAL
}

func GetSyntax(reader *bufio.Reader) SyntaxValue {
	token := GetToken(reader)

	if token.Type == SYMBOL || token.Type == SYMBOL_STRING {
		return SymbolValue{token}
	}

	if token.Type == TOKEN_LPAR {
		return ListValue{getSyntax(reader)}
	} else if token.Type == TOKEN_LIST_LPAR {
		return ListLiteralValue{getSyntax(reader)}
	}

	if token.Type == END {
		return nil
	}

	panic("Syntax error.")
}

func getSyntax(reader *bufio.Reader) []SyntaxValue {
	list := []SyntaxValue{}
	token := GetToken(reader)

	for token.Type != END && token.Type != TOKEN_RPAR {
		if token.Type == SYMBOL || token.Type == SYMBOL_STRING {
			list = append(list, SymbolValue{token})
		} else if token.Type == TOKEN_LPAR {
			list = append(list, ListValue{getSyntax(reader)})
		} else if token.Type == TOKEN_LIST_LPAR {
			list = append(list, ListLiteralValue{getSyntax(reader)})
		} else if token.Type == TOKEN_RPAR || token.Type == TOKEN_LIST_RPAR {
			break
		}

		token = GetToken(reader)
	}

	if token.Type == END {
		panic("Unxpected end of file.")
	}

	return list
}
