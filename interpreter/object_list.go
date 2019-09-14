package interpreter

// List object

type ListObject struct {
	// TODO: implement hashtable
	List []Object
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
	listObject := args[0].(ListObject)
	indexObject := args[1].(NumberObject)

	if len(listObject.List) < indexObject.Integer {
		panic("Index out of range.")
	}

	return listObject.List[indexObject.Integer]
}

func setList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	indexObject := args[1].(NumberObject)
	listObject.List[indexObject.Integer] = args[2]
	return args[2]
}

func appendList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	return ListObject{append(listObject.List, args[1])}
}

func lenList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	return NumberObject{len(listObject.List)}
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
	start := args[1].(NumberObject).Integer
	end := args[2].(NumberObject).Integer

	return ListObject{args[0].(ListObject).List[start:end]}
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
