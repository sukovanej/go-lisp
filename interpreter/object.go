package interpreter

// Object

type Object interface {
	GetSlots() map[string]Object
}

func GetSlot(o Object, slotName string) (Object, bool) {
	slots := o.GetSlots()
	item, ok := slots[slotName]
	return item, ok
}

// Env

type Env struct {
	Objects map[string]Object
	Parent  *Env
}

func (env *Env) GetEnvSymbol(name string) (Object, bool) {
	item, ok := env.Objects[name]
	return item, ok
}

// Number object

type NumberObject struct {
	Integer int
}

func (n NumberObject) GetSlots() map[string]Object {
	return map[string]Object{
		"+": CallableObject{func(args []Object, env *Env) Object {
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
