package interpreter

import "strconv"

// Number object

type NumberObject struct {
	Integer int
}

func addNumbers(args []Object, env *Env) Object {
	first := args[0].(NumberObject)
	second := args[1].(NumberObject)

	return NumberObject{first.Integer + second.Integer}
}

func subtractNumbers(args []Object, env *Env) Object {
	first := args[0].(NumberObject)
	second := args[1].(NumberObject)

	return NumberObject{first.Integer - second.Integer}
}

func multiplyNumbers(args []Object, env *Env) Object {
	first := args[0].(NumberObject)
	second := args[1].(NumberObject)

	return NumberObject{first.Integer * second.Integer}
}

func divideNumbers(args []Object, env *Env) Object {
	first := args[0].(NumberObject)
	second := args[1].(NumberObject)

	return NumberObject{first.Integer / second.Integer}
}

func equalNumbers(args []Object, env *Env) Object {
	first := args[0].(NumberObject)

	switch args[1].(type) {
	case NumberObject:
		second := args[1].(NumberObject)
		return BoolObject{first.Integer == second.Integer}
	default:
		return BoolObject{false}
	}
}

func strNumbers(args []Object, _ *Env) Object {
	value := args[0].(NumberObject)

	return StringObject{strconv.Itoa(value.Integer)}
}

func (n NumberObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__+__":   CallableObject{addNumbers},
		"__-__":   CallableObject{subtractNumbers},
		"__*__":   CallableObject{multiplyNumbers},
		"__/__":   CallableObject{divideNumbers},
		"__==__":  CallableObject{equalNumbers},
		"__str__": CallableObject{strNumbers},
	}
}
