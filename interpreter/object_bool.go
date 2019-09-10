package interpreter

type BoolObject struct {
	Value bool
}

func (obj BoolObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__str__": CallableObject{func(args []Object, _ *Env) Object {
			if args[0].(BoolObject).Value {
				return StringObject{"#t"}
			} else {
				return StringObject{"#f"}
			}
		}},
	}
}
