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

type DictLiteralValue struct {
	Value []SyntaxValue
}

func (_ DictLiteralValue) GetType() SyntaxType {
	return SYNTAX_DICT_LITERAL
}

func GetSyntax(reader *bufio.Reader, meta *BufferMetaInformation) SyntaxValue {
	token := GetToken(reader, meta)

	if token.Type == SYMBOL || token.Type == SYMBOL_STRING {
		return SymbolValue{token}
	}

	if token.Type == TOKEN_LPAR {
		return ListValue{getSyntax(reader, meta)}
	} else if token.Type == TOKEN_LIST_LPAR {
		return ListLiteralValue{getSyntax(reader, meta)}
	} else if token.Type == TOKEN_DICT_LPAR {
		return DictLiteralValue{getSyntax(reader, meta)}
	}

	if token.Type == END {
		return nil
	}

	panic("Syntax error.")
}

func getSyntax(reader *bufio.Reader, meta *BufferMetaInformation) []SyntaxValue {
	list := []SyntaxValue{}
	token := GetToken(reader, meta)

	for token.Type != END && token.Type != TOKEN_RPAR {
		if token.Type == SYMBOL || token.Type == SYMBOL_STRING {
			list = append(list, SymbolValue{token})
		} else if token.Type == TOKEN_LPAR {
			list = append(list, ListValue{getSyntax(reader, meta)})
		} else if token.Type == TOKEN_LIST_LPAR {
			list = append(list, ListLiteralValue{getSyntax(reader, meta)})
		} else if token.Type == TOKEN_DICT_LPAR {
			list = append(list, DictLiteralValue{getSyntax(reader, meta)})
		} else if token.Type == TOKEN_RPAR || token.Type == TOKEN_LIST_RPAR || token.Type == TOKEN_DICT_RPAR {
			break
		}

		token = GetToken(reader, meta)

		if token.Type == END {
			panic("Unxpected end of file.")
		}
	}

	if token.Type == END {
		panic("Unxpected end of file.")
	}

	return list
}
