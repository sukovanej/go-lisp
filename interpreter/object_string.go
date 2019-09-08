package interpreter

// String object

type StringObject struct {
	String string
}

func addStrings(args []Object, env *Env) Object {
	first := args[0].(StringObject)
	second := args[1].(StringObject)

	return StringObject{first.String + second.String}
}

func equalStrings(args []Object, env *Env) Object {
	first := args[0].(StringObject)
	second := args[1].(StringObject)

	return BoolObject{first.String == second.String}
}

func (o StringObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__+__":  CallableObject{addStrings},
		"__==__": CallableObject{equalStrings},
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return o
		}},
	}
}
