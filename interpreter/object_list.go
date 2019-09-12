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

			if !Equal(firstObject, secondObject, env).Value {
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
	listObject.List = append(listObject.List, args[1])
	nilObject, _ := env.GetEnvSymbol("#nil")
	return nilObject
}

func lenList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	return NumberObject{len(listObject.List)}
}

func strList(args []Object, env *Env) Object {
	listObject := args[0].(ListObject)
	l := len(listObject.List)
	result := "["
	for i, entry := range listObject.List {
		result += GetStr(entry, env).String

		if i < l-1 {
			result += ", "
		}
	}
	result += "]"
	return StringObject{result}
}

func (o ListObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__item__":     CallableObject{getList},
		"__len__":      CallableObject{lenList},
		"__set-item__": CallableObject{setList},
		"__==__":       CallableObject{equalLists},
		"__str__":      CallableObject{strList},
		"append":       CallableObject{appendList},
	}
}
