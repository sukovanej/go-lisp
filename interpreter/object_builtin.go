package interpreter

func SetForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	obj := EvalSyntax(args[1], env)
	env.Objects[args[0].(Token).Symbol] = obj
	return obj
}

func AssertCallable(args []Object, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	eqFunc, _ := env.GetEnvSymbol("==")
	if !eqFunc.(CallableObject).Callable(args, env).(BoolObject).Value {
		panic("Assertion error")
	}
	nilObject, _ := env.GetEnvSymbol("#nil")
	return nilObject
}
