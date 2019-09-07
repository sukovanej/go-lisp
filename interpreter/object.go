package interpreter

type Object interface {
	GetSlots() map[string]Object
}

type Env struct {
	Objects map[string]Object
	Parent  *Env
}

type SlotNotFoundError struct {
	error
}

func GetSlot(o Object, slotName string) (Object, bool) {
	slots := o.GetSlots()
	item, ok := slots[slotName]
	return item, ok
}

// Number object

type NumberObject struct {
	Integer int64
}

func (n NumberObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__+__": CallableObject{func(args []Object, env *Env) Object {
			first := args[0].(NumberObject)
			second := args[1].(NumberObject)

			return NumberObject{first.Integer + second.Integer}
		}},
	}
}

// Function object

type CallableObject struct {
	Callable func([]Object, *Env) Object
}

func (n CallableObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__call__": n,
	}
}
