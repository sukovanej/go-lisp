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

	if token.Type == SYMBOL {
		return SyntaxValue(token)
	}

	if token.Type == TOKEN_LPAR {
		list := []SyntaxValue{}
		token = GetToken(reader)

		for token.Type != END {
			if token.Type == SYMBOL {
				list = append(list, SyntaxValue(token))
			} else if token.Type == TOKEN_LPAR {
				list = append(list, GetSyntax(reader))
			} else if token.Type == TOKEN_RPAR {
				break
			}

			token = GetToken(reader)
		}
		return SyntaxValue(list)
	}

	if token.Type == END {
		panic("Unxpected end of file.")
	}

	panic("Syntax error.")
}
