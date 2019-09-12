package interpreter

import "fmt"

func OperatorFunc(operatorName string) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		operatorFunc, ok := GetSlot(args[0], "__"+operatorName+"__")
		if !ok {
			panic("Operator slot not found.")
		}
		operatorCallable := operatorFunc.(CallableObject).Callable
		result := args[0]
		for _, obj := range args[1:len(args)] {
			result = operatorCallable([]Object{result, obj}, env)
		}
		return result
	}
}

func SetForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	obj := EvalSyntax(args[1], env)
	env.SetSymbol(args[0].(Token).Symbol, obj)
	return obj
}

func AssertCallable(args []Object, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	eqFunc, _ := env.GetEnvSymbol("==")
	if !eqFunc.(CallableObject).Callable(args, env).(BoolObject).Value {
		panic("Assertion error")
	}
	nilObject, _ := env.GetEnvSymbol("#nil")
	return nilObject
}

func PrintCallable(args []Object, env *Env) Object {
	for _, obj := range args {
		operatorFunc, ok := GetSlot(obj, "__str__")
		if !ok {
			panic("__str__ slot not found.")
		}
		stringObject := operatorFunc.(CallableObject).Callable([]Object{obj}, env)
		switch stringObject.(type) {
		case StringObject:
			fmt.Print(stringObject.(StringObject).String)
		default:
			panic("__str__ must return string.")
		}
	}

	nilObject, _ := env.GetEnvSymbol("#nil")
	return nilObject
}

func PrintlnCallable(args []Object, env *Env) Object {
	result := PrintCallable(args, env)
	fmt.Print("\n")
	return result
}

func CreateStructForm(slots []SyntaxValue, env *Env) Object {
	slotObjects := map[string]Object{}
	nilObject, _ := env.GetEnvSymbol("#nil")

	for _, slot := range slots {
		slotObjects[slot.(Token).Symbol] = nilObject
	}

	slotObjects["__str__"] = CallableObject{func(obj []Object, env *Env) Object {
		result := "<struct"

		for name, item := range obj[0].GetSlots() {
			result += fmt.Sprintf("%s=%v", name, item)
		}

		result += ">"
		return StringObject{result}
	}}

	return StructObject{slotObjects}
}

func DefStructForm(declared_args []SyntaxValue, env *Env) Object {
	constructor := CallableObject{func(args []Object, _ *Env) Object {
		structObject := CreateStructForm(declared_args[1:], env)
		structSlots := structObject.GetSlots()
		for i, _ := range declared_args[1:] {
			structSlots[declared_args[i].(Token).Symbol] = args[i]
		}
		return structObject
	}}

	env.SetSymbol(declared_args[0].(Token).Symbol, constructor)
	return constructor
}

func GetAttrForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Unexpected number of arguments")
	}
	obj := EvalSyntax(args[0], env)
	slot, ok := GetSlot(obj, args[1].(Token).Symbol)
	if !ok {
		panic("Slot not found")
	}
	return slot
}

func SetAttrForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 3 {
		panic("Unexpected number of arguments")
	}

	slots := EvalSyntax(args[0], env).GetSlots()
	_, ok := slots[args[1].(Token).Symbol]
	if !ok {
		panic(fmt.Sprintf("Struct doesn't have slot '%s'.", args[1].(Token).Symbol))
	}

	object := EvalSyntax(args[2], env)
	slots[args[1].(Token).Symbol] = object
	return object
}

func NotEqualOperator(args []Object, env *Env) Object {
	result := OperatorFunc("==")(args, env).(BoolObject)
	return BoolObject{!result.Value}
}

func StrCallable(args []Object, env *Env) Object {
	operatorFunc, ok := GetSlot(args[0], "__str__")
	if !ok {
		panic("__str__ slot not found.")
	}
	operatorCallable := operatorFunc.(CallableObject).Callable
	return operatorCallable(args, env)
}

func ImportCallable(args []Object, env *Env) Object {
}
