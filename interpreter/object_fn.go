package interpreter

import "fmt"

type FormObject struct {
	Callable func([]SyntaxValue, *Env) Object
}

func (o FormObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__call__": o,
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return StringObject{fmt.Sprintf("<form at %p>", &o)}
		}},
	}
}

func createLambdaFunction(declared_args []SyntaxValue, body []SyntaxValue, env *Env) CallableObject {
	return CallableObject{func(args []Object, _ *Env) Object {
		internal_env := &Env{map[string]Object{}, env}
		isVariadic := declared_args[len(declared_args)-1].(SymbolValue).Value.Type == SYMBOL_VARIADIC

		if len(declared_args) != len(args) && !isVariadic {
			return NewErrorWithSyntaxValue(declared_args[0], fmt.Sprintf("Expected %d arguments, %d given.", len(declared_args), len(args)))
		}

		i := 0
		for _, declared_arg := range declared_args {
			switch declared_arg.GetType() {
			case SYNTAX_SYMBOL:
				if declared_arg.(SymbolValue).Value.Type == SYMBOL_VARIADIC {
					break
				}
				internal_env.Objects[declared_arg.(SymbolValue).Value.Symbol] = args[i]
			default:
				panic("not defined behaviour")
			}
			i++
		}

		if isVariadic {
			internal_env.Objects[declared_args[len(declared_args)-1].(SymbolValue).Value.Symbol] = ListObject{args[i-1:]}
		}

		var last Object
		for _, statement := range body {
			last = EvalSyntax(statement, internal_env)
			if IsErrorObject(last) {
				return last.(ErrorObject).TraceErrorWithSyntaxValue(statement, "")
			}
		}
		return last
	}}
}

func CreateLambdaForm(args []SyntaxValue, env *Env) Object {
	if len(args) < 2 {
		panic("Wrong number of arguments")
	}

	return createLambdaFunction(args[0].(ListValue).Value, args[1:], env)
}

func DefLambdaForm(args []SyntaxValue, env *Env) Object {
	lambda := CreateLambdaForm(args[1:], env)
	env.Objects[args[0].(SymbolValue).Value.Symbol] = lambda
	return lambda
}
