package interpreter

// Dict object

import "fmt"

func IsDictObject(obj Object) bool {
	switch obj.(type) {
	case DictObject:
		return true
	default:
		return false
	}
}

type DictObjectEntry struct {
	Key   Object
	Value Object
}

type DictObject struct {
	Dict map[string]DictObjectEntry
}

func (d DictObject) Len() int {
	return len(d.Dict)
}

func (d DictObject) Get(key Object, env *Env) Object {
	hashObject := GetHash(key, env)
	if IsErrorObject(hashObject) {
		return hashObject
	}
	hash := hashObject.(StringObject).String
	result, ok := d.Dict[hash]
	if !ok {
		panic("not found")
	}
	return result.Value
}

func (d DictObject) Set(key Object, value Object, env *Env) Object {
	nilObject, _ := env.GetEnvSymbol("#nil")
	hashObject := GetHash(key, env)
	if IsErrorObject(hashObject) {
		return hashObject
	}
	hash := hashObject.(StringObject).String
	d.Dict[hash] = DictObjectEntry{key, value}
	return nilObject
}

func GetHash(obj Object, env *Env) Object {
	hashObject, ok := GetSlot(obj, "__hash__")

	if !ok {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be hashable."))
	}

	if IsStringObject(hashObject) {
		return hashObject.(StringObject)
	} else if IsCallableObject(hashObject) {
		hashObject = hashObject.(CallableObject).Callable([]Object{obj}, env)
		if IsStringObject(hashObject) {
			return hashObject.(StringObject)
		}
	}

	return NewErrorWithoutToken(fmt.Sprintf("Hash object must be string or callable which return string."))
}

func equalDicts(args []Object, env *Env) Object {
	first := args[0].(DictObject)

	switch args[1].(type) {
	case DictObject:
		second := args[1].(DictObject)

		if first.Len() != second.Len() {
			return BoolObject{false}
		}

		for _, firstEntry := range first.Dict {
			secondEntry := second.Get(firstEntry.Key, env)
			equal := Equal(firstEntry.Value, secondEntry, env)

			if IsErrorObject(equal) {
				return equal
			}

			if secondEntry == nil || !equal.(BoolObject).Value {
				return BoolObject{false}
			}
		}
	default:
		return BoolObject{false}
	}
	return BoolObject{true}
}

func getDict(args []Object, env *Env) Object {
	if !IsDictObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be dict."))
	}

	dictObject := args[1].(DictObject)
	keyObject := args[0]

	return dictObject.Get(keyObject, env)
}

func setDict(args []Object, env *Env) Object {
	if !IsDictObject(args[2]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be dict."))
	}
	dictObject := args[2].(DictObject)
	keyObject := args[0]
	valueObject := args[1]

	return dictObject.Set(keyObject, valueObject, env)
}

func lenDict(args []Object, env *Env) Object {
	if !IsDictObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be dict."))
	}
	dictObject := args[0].(DictObject)

	return NumberObject{dictObject.Len()}
}

func strDict(args []Object, env *Env) Object {
	dictObject := args[0].(DictObject)
	result := "{"
	for _, entry := range dictObject.Dict {
		keyObject := GetStr(entry.Key, env)
		if IsErrorObject(keyObject) {
			return keyObject
		}
		valueObject := GetStr(entry.Value, env)
		if IsErrorObject(valueObject) {
			return valueObject
		}

		result += keyObject.(StringObject).String + ": " + valueObject.(StringObject).String + ", "
	}
	result = result[:len(result)-2] + "}"
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
