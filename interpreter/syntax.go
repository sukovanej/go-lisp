package interpreter

import (
	"bufio"
)

type SyntaxType int

const (
	SYNTAX_SYMBOL SyntaxType = iota
	SYNTAX_LIST
)

type SyntaxValue interface{}

func GetSyntax(reader *bufio.Reader) SyntaxValue {
	token := GetToken(reader)

	if token.Type == SYMBOL || token.Type == SYMBOL_STRING {
		return SyntaxValue(token)
	}

	if token.Type == TOKEN_LPAR {
		return getSyntax(reader)
	}

	if token.Type == END {
		return nil
	}

	panic("Syntax error.")
}

func getSyntax(reader *bufio.Reader) SyntaxValue {
	list := []SyntaxValue{}
	token := GetToken(reader)

	for token.Type != END && token.Type != TOKEN_RPAR {
		if token.Type == SYMBOL || token.Type == SYMBOL_STRING {
			list = append(list, SyntaxValue(token))
		} else if token.Type == TOKEN_LPAR {
			list = append(list, getSyntax(reader))
		} else if token.Type == TOKEN_RPAR {
			break
		}

		token = GetToken(reader)
	}

	if token.Type == END {
		panic("Unxpected end of file.")
	}

	return SyntaxValue(list)
}
