package interpreter

import (
	"bufio"
	"io"
)

type TokenType int

const (
	TOKEN_LPAR TokenType = iota
	TOKEN_RPAR
	SYMBOL
	END
)

type Token struct {
	Symbol string
	Type   TokenType
}

func GetToken(reader *bufio.Reader) Token {
	r, _, err := reader.ReadRune()

	if r == '(' {
		return Token{string(r), TOKEN_LPAR}
	} else if r == ')' {
		possibleSpace, _, _ := reader.ReadRune()
		if possibleSpace != ' ' {
			reader.UnreadRune()
		}
		return Token{string(r), TOKEN_RPAR}
	} else if err == io.EOF {
		return Token{"", END}
	}

	result := string(r)

	for {
		r, _, err = reader.ReadRune()

		if r == ' ' {
			break
		}

		if err == io.EOF || r == ')' {
			reader.UnreadRune()
			break
		}
		result += string(r)
	}

	return Token{result, SYMBOL}
}
