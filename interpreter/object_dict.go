package interpreter

// Dict object

type DictObject struct {
	dict map[string]Object
}

func (d DictObject) Get(key Object, env *Env) Object {
	hash := GetHash(key, env)
	obj, ok = d.dict[hash]

	if !ok {
		panic("Item not found")
	}

	return obj
}

func (d DictObject) Set(key Object, value Object, env *Env) Object {
	nilObject := env.GetEnvSymbol("#nil")
	hash := GetHash(key, env)
	d.dict[hash] = value

	if !ok {
		panic("Item not found")
	}

	return nilObject
}

func GetHash(obj Object, env *Env) string {
	hashObject, ok := GetSlot(obj, "__hash__")

	if !ok {
		panic("Object is not hashable.")
	}

	switch hashObject.(type) {
	case StringObject:
		return hashObject.(StringObject).Value
	case CallableObject:
		hashObject = hashObject.(CallableObject).Callable([]Object{obj}, env)
		switch hashObject.(type) {
		case StringObject:
			return hashObject.(StringObject).Value
		}
	}

	panic("__hash__ must be str or callable")
}

func equalDicts(args []Object, env *Env) Object {
	first := args[0].(DictObject)

	switch args[1].(type) {
	case DictObject:
		second := args[1].(DictObject)

		return BoolObject{first.String == second.String}
	default:
		return BoolObject{false}
	}
}

func (o DictObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__get__": CallableObject{addStrings},
		"__set__": CallableObject{equalStrings},
		"__==__":  CallableObject{equalDicts},
	}
}
