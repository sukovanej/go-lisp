package interpreter

import "fmt"

type FormObject struct {
	Callable func([]Object, *Env) Object
}

func (o FormObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__call__": o,
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return StringObject{fmt.Sprintf("<form at %p>", &o)}
		}},
	}
}

func createCallable(declared_args []SyntaxValue, body []SyntaxValue, env *Env) func([]Object, *Env) Object {
	return func(args []Object, _ *Env) Object {
		internal_env := &Env{map[string]Object{}, env}
		isVariadic := len(declared_args) != 0 && declared_args[len(declared_args)-1].(SymbolValue).Value.Type == SYMBOL_VARIADIC

		if len(declared_args) != len(args) && !isVariadic {
			// TODO: error when declared_args is empty
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
	}
}

func createLambdaFunction(declared_args []SyntaxValue, body []SyntaxValue, env *Env) CallableObject {
	return CallableObject{createCallable(declared_args, body, env)}
}

func createLambdaForm(declared_args []SyntaxValue, body []SyntaxValue, env *Env) FormObject {
	return FormObject{createCallable(declared_args, body, env)}
}

func CreateLambdaForm(args []Object, env *Env) Object {
	if len(args) < 2 {
		panic("Wrong number of arguments")
	}

	bodyArgs := []SyntaxValue{}
	for _, arg := range args[1:] {
		bodyArgs = append(bodyArgs, arg.(SyntaxObject).Value)
	}

	return createLambdaFunction(args[0].(SyntaxObject).Value.(ListValue).Value, bodyArgs, env)
}

func DefLambdaForm(args []Object, env *Env) Object {
	lambda := CreateLambdaForm(args[1:], env)
	env.Objects[args[0].(SyntaxObject).Value.(SymbolValue).Value.Symbol] = lambda
	return lambda
}

func IsCallableObject(obj Object) bool {
	switch obj.(type) {
	case CallableObject:
		return true
	default:
		return false
	}
}

func CreateFormForm(args []Object, env *Env) Object {
	if len(args) < 2 {
		return NewErrorWithSyntaxValue(args[0].(SyntaxObject).Value, fmt.Sprintf("Expected 2 arguments, %d given.", len(args)))
	}

	bodyArgs := []SyntaxValue{}
	for _, arg := range args[1:] {
		bodyArgs = append(bodyArgs, arg.(SyntaxObject).Value)
	}

	return createLambdaForm(args[0].(SyntaxObject).Value.(ListValue).Value, bodyArgs, env)
}

func DefFormForm(args []Object, env *Env) Object {
	lambda := CreateFormForm(args[1:], env)
	env.Objects[args[0].(SyntaxObject).Value.(SymbolValue).Value.Symbol] = lambda
	return lambda
}
