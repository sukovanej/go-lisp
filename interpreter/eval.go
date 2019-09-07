package interpreter

import (
	"bufio"
)

func Eval(reader *bufio.Reader, env *Env) Object {
	syntaxTree := GetSyntax(reader)
	return eval(syntaxTree, env)
}

func eval(value SyntaxValue, env *Env) Object {
	switch value.(type) {
	case Token:
		return evalSymbol(value.(Token), env)
	case []Token:
		return evalFunction(value.([]SyntaxValue), env)
	}
	return nil
}

func evalSymbol(token Token, env *Env) Object {
	return nil
}

func evalFunction(list []SyntaxValue, env *Env) Object {
	function := eval(list[0], env)
	args := []Object{}

	for arg, _ := range list {
		args = append(args, eval(arg, env))
	}

	return function.(CallableObject).Callable(args, env)
}

func GetMainEnv() *Env {
	return &Env{
		map[string]Object{
			"+": CallableObject{func(args []Object, env *Env) Object {
				plusFunc, ok := GetSlot(args[0], "__+__")
				if !ok {
					panic("Symbol not found")
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
