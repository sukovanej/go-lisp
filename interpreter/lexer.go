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

type BufferMetaInformation struct {
	Line     int
	Column   int
	Filename string
}

func (b *BufferMetaInformation) Incr(r rune) {
	if r == '\n' {
		b.Column = 0
		b.Line++
	} else {
		b.Column++
	}
}

type Token struct {
	Symbol   string
	Type     TokenType
	Line     int
	Column   int
	Filename string
}

func skipCommentsAndWhitespaces(r rune, err error, reader *bufio.Reader, meta *BufferMetaInformation) (rune, error) {
	// skip comments
	if r == ';' {
		for r != '\n' && err != io.EOF {
			r, _, err = reader.ReadRune()
			meta.Incr(r)
		}
	}

	for r == ' ' || r == '\n' || r == '\t' {
		r, _, err = reader.ReadRune()
		meta.Incr(r)
	}

	if r == ';' {
		return skipCommentsAndWhitespaces(r, err, reader, meta)
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

func GetToken(reader *bufio.Reader, meta *BufferMetaInformation) Token {
	r, _, err := reader.ReadRune()
	meta.Incr(r)
	r, err = skipCommentsAndWhitespaces(r, err, reader, meta)
	line, column, file := meta.Line, meta.Column, meta.Filename

	if r == '(' {
		return Token{string(r), TOKEN_LPAR, line, column, file}
	} else if r == '[' {
		return Token{string(r), TOKEN_LIST_LPAR, line, column, file}
	} else if r == '{' {
		return Token{string(r), TOKEN_DICT_LPAR, line, column, file}
	} else if r == ')' {
		return Token{string(r), TOKEN_RPAR, line, column, file}
	} else if r == ']' {
		return Token{string(r), TOKEN_LIST_RPAR, line, column, file}
	} else if r == '}' {
		return Token{string(r), TOKEN_DICT_RPAR, line, column, file}
	} else if r == '"' {
		str := ""
		r, _, err = reader.ReadRune()
		meta.Incr(r)
		for r != '"' {
			if r == '\\' {
				r, _, err = reader.ReadRune()
				meta.Incr(r)
				r = createSpecialCharacter(r)
			}
			str += string(r)
			r, _, err = reader.ReadRune()
			meta.Incr(r)
		}
		return Token{str, SYMBOL_STRING, line, column, file}
	} else if err == io.EOF {
		return Token{"", END, line, column, file}
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
		meta.Incr(r)
		result += string(r)
	}

	return Token{result, SYMBOL, line, column, file}
}
