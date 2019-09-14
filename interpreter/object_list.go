package interpreter

// List object

import "fmt"

type ListObject struct {
	// TODO: implement hashtable
	List []Object
}

func IsListObject(obj Object) bool {
	switch obj.(type) {
	case ListObject:
		return true
	default:
		return false
	}
}

func equalLists(args []Object, env *Env) Object {
	first := args[0].(ListObject)

	switch args[1].(type) {
	case ListObject:
		second := args[1].(ListObject)

		if len(first.List) != len(second.List) {
			return BoolObject{false}
		}

		for index, firstObject := range first.List {
			secondObject := second.List[index]
			equal := Equal(firstObject, secondObject, env)

			if IsErrorObject(equal) {
				return equal
			}

			if !equal.(BoolObject).Value {
				return BoolObject{false}
			}
		}
	default:
		return BoolObject{false}
	}
	return BoolObject{true}
}

func getList(args []Object, env *Env) Object {
	if !IsNumberObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be number."))
	} else if !IsListObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be list."))
	}

	index := args[0].(NumberObject).Integer
	list := args[1].(ListObject).List

	if len(list) <= index {
		return NewErrorWithoutToken(fmt.Sprintf("Index out of range."))
	}

	return list[index]
}

func setList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	indexObject := args[1].(NumberObject)
	listObject.List[indexObject.Integer] = args[2]
	return args[2]
}

func appendList(args []Object, env *Env) Object {
	if !IsListObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be list."))
	}
	return ListObject{append(args[1].(ListObject).List, args[0])}
}

func lenList(args []Object, env *Env) Object {
	if !IsListObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be list."))
	}
	list := args[0].(ListObject).List
	return NumberObject{len(list)}
}

func plusList(args []Object, env *Env) Object {
	listObjectLeft := args[0].(ListObject)
	listObjectRight := args[1].(ListObject)
	return ListObject{append(listObjectLeft.List, listObjectRight.List...)}
}

func strList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	l := len(listObject.List)
	result := "["
	for i, entry := range listObject.List {
		str := GetStr(entry, env)
		if IsErrorObject(str) {
			return str
		}
		result += str.(StringObject).String

		if i < l-1 {
			result += ", "
		}
	}
	result += "]"
	return StringObject{result}
}

func sliceList(args []Object, env *Env) Object {
	if !IsNumberObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be number."))
	} else if !IsNumberObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be number."))
	} else if !IsListObject(args[2]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be list."))
	}

	list := args[2].(ListObject).List
	left := args[0].(NumberObject).Integer
	right := args[1].(NumberObject).Integer

	return ListObject{list[left:right]}
}

func (o ListObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__item__":     CallableObject{getList},
		"__len__":      CallableObject{lenList},
		"__set-item__": CallableObject{setList},
		"__==__":       CallableObject{equalLists},
		"__str__":      CallableObject{strList},
		"__+__":        CallableObject{plusList},
		"__slice__":    CallableObject{sliceList},
		"__append__":   CallableObject{appendList},
	}
}
