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
		if len(declared_args) != len(args) {
			panic("wrong number of arguments")
		}
		for i, declared_arg := range declared_args {
			switch declared_arg.GetType() {
			case SYNTAX_SYMBOL:
				internal_env.Objects[declared_arg.(SymbolValue).Value.Symbol] = args[i]
			default:
				panic("not defined behaviour")
			}
		}

		var last Object
		for _, statement := range body {
			last = EvalSyntax(statement, internal_env)
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
