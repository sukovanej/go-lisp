package interpreter

func IfForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 3 && len(args) != 2 {
		panic("Wrong number of arguments")
	}

	condition := EvalSyntax(args[0], env)

	if IsErrorObject(condition) {
		return condition
	}

	switch condition.(type) {
	case BoolObject:
		if condition.(BoolObject).Value {
			result := EvalSyntax(args[1], env)
			if IsErrorObject(result) {
				return result
			}
			return result
		}

		if len(args) == 3 {
			result := EvalSyntax(args[2], env)
			if IsErrorObject(result) {
				return result
			}
			return result
		} else {
			n, _ := env.GetEnvSymbol("#nil")
			return n
		}
	default:
		panic("Condition is not bool.")
	}
}
