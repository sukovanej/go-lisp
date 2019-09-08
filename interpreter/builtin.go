package interpreter

func BinaryOperatorFunc(operatorName string) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		plusFunc, ok := GetSlot(args[0], "__"+operatorName+"__")
		if !ok {
			panic("Operator slot not found.")
		}
		plusCallable := plusFunc.(CallableObject).Callable
		result := args[0]
		for _, obj := range args[1:len(args)] {
			result = plusCallable([]Object{result, obj}, nil)
		}
		return result
	}
}
