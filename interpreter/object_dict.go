package interpreter

// Dict object

import "fmt"

type DictObjectEntry struct {
	HashKey string
	Key     Object
	Value   Object
	Next    *DictObjectEntry
}

type DictObject struct {
	// TODO: implement hashtable
	Dict []*DictObjectEntry
}

func (d DictObject) Len() int {
	result := 0
	entry := d.Dict[0]
	for entry != nil {
		result++
		entry = entry.Next
	}
	return result
}

func (d DictObject) Get(key Object, env *Env) Object {
	hash := GetHash(key, env)
	entry := d.Dict[0]

	for entry != nil {
		if hash == entry.HashKey {
			return entry.Value
		}

		entry = entry.Next
	}

	return nil
}

func (d DictObject) Set(key Object, value Object, env *Env) Object {
	nilObject, _ := env.GetEnvSymbol("#nil")
	hash := GetHash(key, env)
	first := d.Dict[0]
	d.Dict[0] = &DictObjectEntry{hash, key, value, first}
	return nilObject
}

func GetHash(obj Object, env *Env) string {
	hashObject, ok := GetSlot(obj, "__hash__")

	if !ok {
		NewErrorWithoutToken(fmt.Sprintf(""))
	}

	switch hashObject.(type) {
	case StringObject:
		return hashObject.(StringObject).String
	case CallableObject:
		hashObject = hashObject.(CallableObject).Callable([]Object{obj}, env)
		switch hashObject.(type) {
		case StringObject:
			return hashObject.(StringObject).String
		}
	}

	panic("__hash__ must be str or callable")
}

func equalDicts(args []Object, env *Env) Object {
	first := args[0].(DictObject)

	switch args[1].(type) {
	case DictObject:
		second := args[1].(DictObject)

		if first.Len() != second.Len() {
			return BoolObject{false}
		}

		firstEntry := first.Dict[0]
		for firstEntry != nil {
			secondEntry := second.Get(firstEntry.Key, env)
			equal := Equal(firstEntry.Value, secondEntry, env)

			if IsErrorObject(equal) {
				return equal
			}

			if secondEntry == nil || !equal.(BoolObject).Value {
				return BoolObject{false}
			}

			firstEntry = firstEntry.Next
		}
	default:
		return BoolObject{false}
	}
	return BoolObject{true}
}

func getDict(args []Object, env *Env) Object {
	dictObject := args[0].(DictObject)
	keyObject := args[1]

	return dictObject.Get(keyObject, env)
}

func setDict(args []Object, env *Env) Object {
	dictObject := args[0].(DictObject)
	keyObject := args[1]
	valueObject := args[2]

	return dictObject.Set(keyObject, valueObject, env)
}

func lenDict(args []Object, env *Env) Object {
	dictObject := args[0].(DictObject)

	return NumberObject{dictObject.Len()}
}

func strDict(args []Object, env *Env) Object {
	dictObject := args[0].(DictObject)
	result := "{"
	entry := dictObject.Dict[0]
	for entry != nil {
		keyObject := GetStr(entry.Key, env)
		if IsErrorObject(keyObject) {
			return keyObject
		}
		valueObject := GetStr(entry.Value, env)
		if IsErrorObject(valueObject) {
			return valueObject
		}

		result += keyObject.(StringObject).String + ": " + valueObject.(StringObject).String
		entry = entry.Next

		if entry != nil {
			result += ", "
		}
	}
	result += "}"
	return StringObject{result}
}

func (o DictObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__item__":     CallableObject{getDict},
		"__len__":      CallableObject{lenDict},
		"__set-item__": CallableObject{setDict},
		"__==__":       CallableObject{equalDicts},
		"__str__":      CallableObject{strDict},
	}
}
