package interpreter

// String object

import "fmt"

type StringObject struct {
	String string
}

func IsStringObject(obj Object) bool {
	switch obj.(type) {
	case StringObject:
		return true
	default:
		return false
	}
}

func addStrings(args []Object, env *Env) Object {
	if !IsStringObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	} else if !IsStringObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	}
	first := args[0].(StringObject).String
	second := args[1].(StringObject).String

	return StringObject{first + second}
}

func equalStrings(args []Object, env *Env) Object {
	if !IsStringObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	} else if !IsStringObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	}

	first := args[0].(StringObject).String
	second := args[1].(StringObject).String

	return BoolObject{first == second}
}

func itemStrings(args []Object, env *Env) Object {
	if !IsNumberObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be number."))
	} else if !IsStringObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	}

	index := args[0].(NumberObject).Integer
	first := args[1].(StringObject).String

	if len(first) <= index {
		iStringObject := GetStr(args[1], env)
		return NewErrorWithoutToken(fmt.Sprintf("Index %s is out of range.", iStringObject.(StringObject).String))
	}

	return StringObject{string(first[index])}
}

func sliceString(args []Object, env *Env) Object {
	if !IsNumberObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be number."))
	} else if !IsNumberObject(args[1]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be number."))
	} else if !IsStringObject(args[2]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	}

	start := args[0].(NumberObject).Integer
	end := args[1].(NumberObject).Integer
	str := args[2].(StringObject).String

	return StringObject{str[start:end]}
}

func lenString(args []Object, env *Env) Object {
	str := args[0].(StringObject).String

	return NumberObject{len(str)}
}

func appendString(args []Object, env *Env) Object {
	first := args[0].(StringObject).String
	second := args[1].(StringObject).String

	return StringObject{second + first}
}

func (o StringObject) GetSlots() map[string]Object {
	return map[string]Object{
		"__+__":      CallableObject{addStrings},
		"__==__":     CallableObject{equalStrings},
		"__hash__":   StringObject{"__str__" + o.String},
		"__item__":   CallableObject{itemStrings},
		"__slice__":  CallableObject{sliceString},
		"__append__": CallableObject{appendString},
		"__len__":    CallableObject{lenString},
		"__str__": CallableObject{func(_ []Object, _ *Env) Object {
			return o
		}},
	}
}
