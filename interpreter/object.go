package interpreter

import (
	"fmt"
)

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
	if !ok && env.Parent != nil {
		return env.Parent.GetEnvSymbol(name)
	}
	return item, ok
}

func (env *Env) SetSymbol(name string, obj Object) Object {
	env.Objects[name] = obj
	return obj
}

// Nil Object

type NilObject struct{}

func (_ NilObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return StringObject{"Nil"}
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
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return StringObject{fmt.Sprintf("<callable at %p>", &n)}
		}},
	}
}
