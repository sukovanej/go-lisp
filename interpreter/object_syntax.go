package interpreter

import "fmt"

func IsSyntaxObject(obj Object) bool {
	switch obj.(type) {
	case SyntaxObject:
		return true
	default:
		return false
	}
}

type SyntaxObject struct {
	Value SyntaxValue
}

func (s SyntaxObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return StringObject{fmt.Sprintf("<form at %p>", &s)}
		}},
	}
}
