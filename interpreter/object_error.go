package interpreter

import "fmt"

type ErrorObjectEntry struct {
	Token       Token
	Description string
}

type ErrorObject struct {
	Stacktrace []ErrorObjectEntry
}

func getTokenFromSyntaxValue(value SyntaxValue) Token {
	switch value.GetType() {
	case SYNTAX_SYMBOL:
		return value.(SymbolValue).Value
	case SYNTAX_LIST:
		return getTokenFromSyntaxValue(value.(ListValue).Value[0])
	case SYNTAX_LIST_LITERAL:
		return getTokenFromSyntaxValue(value.(ListLiteralValue).Value[0])
	case SYNTAX_DICT_LITERAL:
		return getTokenFromSyntaxValue(value.(DictLiteralValue).Value[0])
	default:
		panic("Unknown syntax token type")
	}
}

func NewErrorWithSyntaxValue(value SyntaxValue, description string) ErrorObject {
	return ErrorObject{[]ErrorObjectEntry{ErrorObjectEntry{getTokenFromSyntaxValue(value), description}}}
}

func NewErrorWithoutToken(description string) ErrorObject {
	return ErrorObject{[]ErrorObjectEntry{ErrorObjectEntry{Token{"", END, -1, -1, ""}, description}}}
}

func NewError(token Token, description string) ErrorObject {
	return ErrorObject{[]ErrorObjectEntry{ErrorObjectEntry{token, description}}}
}

func (e ErrorObject) TraceError(token Token, description string) ErrorObject {
	return ErrorObject{append([]ErrorObjectEntry{ErrorObjectEntry{token, description}}, e.Stacktrace...)}
}

func (e ErrorObject) TraceErrorWithSyntaxValue(value SyntaxValue, description string) ErrorObject {
	return ErrorObject{append([]ErrorObjectEntry{ErrorObjectEntry{getTokenFromSyntaxValue(value), description}}, e.Stacktrace...)}
}

func IsErrorObject(e Object) bool {
	switch e.(type) {
	case ErrorObject:
		return true
	default:
		return false
	}
}

func strError(args []Object, _ *Env) Object {
	return StringObject{"<Error object>"}
}

func PrintTraceback(errorObject ErrorObject) {
	fmt.Println("Error traceback:")
	for _, entry := range errorObject.Stacktrace {
		if entry.Token.Type != END {
			fmt.Printf("In file %s\nnear token %s on %d:%d.\n", entry.Token.Filename, entry.Token.Symbol, entry.Token.Line, entry.Token.Column)
		}
		fmt.Println(entry.Description)
	}
}

func (n ErrorObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__str__": CallableObject{strError},
	}
}
