package interpreter

import (
	"bufio"
	"io"
)

type TokenType int

const (
	TOKEN_LPAR TokenType = iota
	TOKEN_RPAR
	TOKEN_LIST_LPAR
	TOKEN_LIST_RPAR
	TOKEN_DICT_LPAR
	TOKEN_DICT_RPAR
	SYMBOL
	SYMBOL_STRING
	END
)

type Token struct {
	Symbol string
	Type   TokenType
}

func skipCommentsAndWhitespaces(r rune, err error, reader *bufio.Reader) (rune, error) {
	// skip comments
	if r == ';' {
		for r != '\n' && err != io.EOF {
			r, _, err = reader.ReadRune()
		}
	}

	for r == ' ' || r == '\n' || r == '\t' {
		r, _, err = reader.ReadRune()
	}

	if r == ';' {
		return skipCommentsAndWhitespaces(r, err, reader)
	}

	return r, err
}

func createSpecialCharacter(r rune) rune {
	switch r {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case '"':
		return '"'
	default:
		return r
	}
}

func GetToken(reader *bufio.Reader) Token {
	r, _, err := reader.ReadRune()
	r, err = skipCommentsAndWhitespaces(r, err, reader)

	if r == '(' {
		return Token{string(r), TOKEN_LPAR}
	} else if r == '[' {
		return Token{string(r), TOKEN_LIST_LPAR}
	} else if r == '{' {
		return Token{string(r), TOKEN_DICT_LPAR}
	} else if r == ')' {
		return Token{string(r), TOKEN_RPAR}
	} else if r == ']' {
		return Token{string(r), TOKEN_LIST_RPAR}
	} else if r == '}' {
		return Token{string(r), TOKEN_DICT_RPAR}
	} else if r == '"' {
		str := ""
		r, _, err = reader.ReadRune()
		for r != '"' {
			if r == '\\' {
				r, _, err = reader.ReadRune()
				r = createSpecialCharacter(r)
			}
			str += string(r)
			r, _, err = reader.ReadRune()
		}
		return Token{str, SYMBOL_STRING}
	} else if err == io.EOF {
		return Token{"", END}
	}

	result := string(r)

	for {
		r, _, err = reader.ReadRune()

		if r == ' ' || r == '\n' || r == '\t' {
			break
		}

		if err == io.EOF || r == ')' || r == '}' || r == ']' {
			reader.UnreadRune()
			break
		}
		result += string(r)
	}

	return Token{result, SYMBOL}
}
