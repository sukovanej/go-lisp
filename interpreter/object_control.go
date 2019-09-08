package interpreter

func IfForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 3 && len(args) != 2 {
		panic("Wrong number of arguments")
	}

	condition := EvalSyntax(args[0], env)
	switch condition.(type) {
	case BoolObject:
		if condition.(BoolObject).Value {
			return EvalSyntax(args[1], env)
		}

		if len(args) == 3 {
			return EvalSyntax(args[2], env)
		} else {
			n, _ := env.GetEnvSymbol("#nil")
			return n
		}
	default:
		panic("Condition is not bool.")
	}
}
