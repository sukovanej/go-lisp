package interpreter

type BoolObject struct {
	Value bool
}

func (obj BoolObject) GetSlots() map[string]Object {
	return map[string]Object{}
}
