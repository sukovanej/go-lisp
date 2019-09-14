package interpreter

type BoolObject struct {
	Value bool
}

func equalBool(args []Object, env *Env) Object {
	left := args[0].(BoolObject).Value
	right := args[1].(BoolObject).Value
	return BoolObject{left == right}
}

func IsBoolObject(obj Object) bool {
	switch obj.(type) {
	case BoolObject:
		return true
	default:
		return false
	}
}

func (obj BoolObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__==__": CallableObject{equalBool},
		"__str__": CallableObject{func(args []Object, _ *Env) Object {
			if args[0].(BoolObject).Value {
				return StringObject{"#t"}
			} else {
				return StringObject{"#f"}
			}
		}},
	}
}
