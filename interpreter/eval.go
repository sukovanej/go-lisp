package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func EvalFile(file string, env *Env) Object {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	return Eval(bufio.NewReader(f), env)
}

func Eval(reader *bufio.Reader, env *Env) Object {
	syntaxTree := GetSyntax(reader)
	var lastResult Object

	for syntaxTree != nil {
		lastResult = EvalSyntax(syntaxTree, env)
		syntaxTree = GetSyntax(reader)
	}

	return lastResult
}

func EvalSyntax(value SyntaxValue, env *Env) Object {
	switch value.(type) {
	case Token:
		return evalSymbol(value.(Token), env)
	case []SyntaxValue:
		return evalFunction(value.([]SyntaxValue), env)
	}

	panic("Unexpected syntax token.")
}

func evalSymbol(token Token, env *Env) Object {
	if token.Type == SYMBOL_STRING {
		return StringObject{token.Symbol}
	} else if num, err := strconv.Atoi(token.Symbol); err == nil {
		return NumberObject{num}
	} else {
		obj, ok := env.GetEnvSymbol(token.Symbol)
		if !ok {
			panic(fmt.Sprintf("Symbol %s not found.", token.Symbol))
		}
		return obj
	}
}

func evalFunction(list []SyntaxValue, env *Env) Object {
	function := EvalSyntax(list[0], env)
	switch function.(type) {
	case CallableObject:
		args := []Object{}

		for _, arg := range list[1:len(list)] {
			value := EvalSyntax(arg, env)
			args = append(args, value)
		}

		return function.(CallableObject).Callable(args, env)
	case FormObject:
		args := list[1:len(list)]
		return function.(FormObject).Callable(args, env)
	default:
		panic("First item must be callable")
	}
}

func GetMainEnv() *Env {
	return &Env{
		map[string]Object{
			"+":       CallableObject{BinaryOperatorFunc("+")},
			"-":       CallableObject{BinaryOperatorFunc("-")},
			"*":       CallableObject{BinaryOperatorFunc("*")},
			"/":       CallableObject{BinaryOperatorFunc("/")},
			"^":       CallableObject{BinaryOperatorFunc("^")},
			"fn":      FormObject{CreateLambdaForm},
			"set":     FormObject{SetForm},
			"#t":      BoolObject{true},
			"#f":      BoolObject{false},
			"==":      CallableObject{BinaryOperatorFunc("==")},
			"#nil":    NilObject{},
			"if":      FormObject{IfForm},
			"!assert": CallableObject{AssertCallable},
			"print":   CallableObject{PrintCallable},
		},
		nil,
	}

}
