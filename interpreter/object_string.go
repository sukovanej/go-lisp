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

	switch args[1].(type) {
	case StringObject:
		second := args[1].(StringObject)

		return BoolObject{first.String == second.String}
	default:
		return BoolObject{false}
	}
}

func itemStrings(args []Object, env *Env) Object {
	first := args[0].(StringObject).String
	i := args[1].(NumberObject).Integer
	return StringObject{string(first[i])}
}

func sliceString(args []Object, env *Env) Object {
	str := args[0].(StringObject).String
	start := args[1].(NumberObject).Integer
	end := args[2].(NumberObject).Integer

	return StringObject{str[start:end]}
}

func lenString(args []Object, env *Env) Object {
	str := args[0].(StringObject).String

	return NumberObject{len(str)}
}

func (o StringObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__+__":     CallableObject{addStrings},
		"__==__":    CallableObject{equalStrings},
		"__hash__":  StringObject{"__str__" + o.String},
		"__item__":  CallableObject{itemStrings},
		"__slice__": CallableObject{sliceString},
		"__len__":   CallableObject{lenString},
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return o
		}},
	}
}
