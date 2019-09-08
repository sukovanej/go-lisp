package interpreter

import (
	"bufio"
	"fmt"
	"strconv"
)

func Eval(reader *bufio.Reader, env *Env) Object {
	syntaxTree := GetSyntax(reader)
	return eval(syntaxTree, env)
}

func eval(value SyntaxValue, env *Env) Object {
	switch value.(type) {
	case Token:
		return evalSymbol(value.(Token), env)
	case []SyntaxValue:
		return evalFunction(value.([]SyntaxValue), env)
	}

	panic("Unexpected syntax token.")
}

func evalSymbol(token Token, env *Env) Object {
	if num, err := strconv.Atoi(token.Symbol); err == nil {
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
	function := eval(list[0], env)
	args := []Object{}

	for _, arg := range list[1:len(list)] {
		value := eval(arg, env)
		args = append(args, value)
	}

	return function.(CallableObject).Callable(args, env)
}

func GetMainEnv() *Env {
	return &Env{
		map[string]Object{
			"+": CallableObject{func(args []Object, env *Env) Object {
				plusFunc, ok := GetSlot(args[0], "+")
				if !ok {
					panic("Object is not addable.")
				}
				plusCallable := plusFunc.(CallableObject).Callable
				result := args[0]
				for _, obj := range args[1:len(args)] {
					result = plusCallable([]Object{result, obj}, nil)
				}
				return result
			}},
		},
		nil,
	}
}
