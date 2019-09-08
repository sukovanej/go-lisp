package interpreter

import "fmt"

func BinaryOperatorFunc(operatorName string) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		operatorFunc, ok := GetSlot(args[0], "__"+operatorName+"__")
		if !ok {
			panic("Operator slot not found.")
		}
		operatorCallable := operatorFunc.(CallableObject).Callable
		result := args[0]
		for _, obj := range args[1:len(args)] {
			result = operatorCallable([]Object{result, obj}, nil)
		}
		return result
	}
}

func SetForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	obj := EvalSyntax(args[1], env)
	env.Objects[args[0].(Token).Symbol] = obj
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
		stringObject := operatorFunc.(CallableObject).Callable([]Object{obj}, nil).(StringObject)
		fmt.Print(stringObject.String)
	}

	nilObject, _ := env.GetEnvSymbol("nil")
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
	return StructObject{slotObjects}
}

func GetAttrForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Unexpected number of arguments")
	}
	slot, ok := GetSlot(EvalSyntax(args[0], env), args[1].(Token).Symbol)
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
