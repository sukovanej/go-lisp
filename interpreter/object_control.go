package interpreter

func IfForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 3 && len(args) != 2 {
		panic("Wrong number of arguments")
	}

	condition := EvalSyntax(args[0], env)

	if IsErrorObject(condition) {
		return condition
	}

	if !IsBoolObject(condition) {
		return NewErrorWithSyntaxValue(args[0], "Condition must be bool object.")
	}

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
	}
	n, _ := env.GetEnvSymbol("#nil")
	return n
}

func CondForm(args []SyntaxValue, env *Env) Object {
	if len(args)%2 != 0 || len(args) < 2 {
		return NewErrorWithSyntaxValue(args[0], "Cond form expects arguments of form (cond condition-1 statement-1 condition-2 statement-2 ...)")
	}

	for i := 0; i < len(args); i += 2 {
		conditionObject := EvalSyntax(args[i], env)
		if IsErrorObject(conditionObject) {
			return conditionObject
		}

		if !IsBoolObject(conditionObject) {
			return NewErrorWithSyntaxValue(args[i], "Condition must be bool object.")
		}

		if conditionObject.(BoolObject).Value {
			return EvalSyntax(args[i+1], env)
		}
	}

	n, _ := env.GetEnvSymbol("#nil")
	return n
}
