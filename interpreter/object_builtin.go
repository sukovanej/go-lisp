package interpreter

import "bufio"
import "fmt"
import "os"

func OperatorFunc(operatorName string) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		operatorFunc, ok := GetSlot(args[0], "__"+operatorName+"__")
		if !ok {
			panic(fmt.Sprintf("Operator %s slot not found on %v.", operatorName, args[0]))
		}
		operatorCallable := operatorFunc.(CallableObject).Callable
		result := args[0]
		for _, obj := range args[1:len(args)] {
			result = operatorCallable([]Object{result, obj}, env)
		}
		return result
	}
}

func Equal(left Object, right Object, env *Env) BoolObject {
	operatorFunc, ok := GetSlot(left, "__==__")
	if !ok {
		panic("__==__ slot not found.")
	}
	operatorCallable := operatorFunc.(CallableObject).Callable
	return operatorCallable([]Object{left, right}, env).(BoolObject)
}

func SetForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	obj := EvalSyntax(args[1], env)
	env.SetSymbol(args[0].(SymbolValue).Value.Symbol, obj)
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
		slotObjects[slot.(SymbolValue).Value.Symbol] = nilObject
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
			structSlots[declared_args[i+1].(SymbolValue).Value.Symbol] = args[i]
		}
		return structObject
	}}

	env.SetSymbol(declared_args[0].(SymbolValue).Value.Symbol, constructor)
	return constructor
}

func GetAttrForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Unexpected number of arguments")
	}
	obj := EvalSyntax(args[0], env)
	slot, ok := GetSlot(obj, args[1].(SymbolValue).Value.Symbol)
	if !ok {
		panic(fmt.Sprintf("Slot %s not found", args[1].(SymbolValue).Value.Symbol))
	}
	return slot
}

func SetAttrForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 3 {
		panic("Unexpected number of arguments")
	}

	slots := EvalSyntax(args[0], env).GetSlots()
	_, ok := slots[args[1].(SymbolValue).Value.Symbol]
	if !ok {
		panic(fmt.Sprintf("Struct doesn't have slot '%s'.", args[1].(SymbolValue).Value.Symbol))
	}

	object := EvalSyntax(args[2], env)
	slots[args[1].(SymbolValue).Value.Symbol] = object
	return object
}

func NotEqualOperator(args []Object, env *Env) Object {
	result := OperatorFunc("==")(args, env).(BoolObject)
	return BoolObject{!result.Value}
}

func SlotCallable(slotName string, numberOfArguments int) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		if len(args) != numberOfArguments {
			panic(fmt.Sprintf("Expected %d arguments, %d given.", numberOfArguments, len(args)))
		}
		operatorFunc, ok := GetSlot(args[0], slotName)
		if !ok {
			panic(fmt.Sprintf("%s slot not found.", slotName))
		}
		operatorCallable := operatorFunc.(CallableObject).Callable
		return operatorCallable(args, env)
	}
}

func ImportCallable(args []Object, env *Env) Object {
	importFileString := args[0].(StringObject).String
	importPathObject, ok := env.GetEnvSymbol("__import_path__")

	if !ok {
		panic("__import_path__ not found")
	}

	for _, pathObject := range importPathObject.(ListObject).List {
		fileName := pathObject.(StringObject).String + "/" + importFileString + ".gisp"
		if _, err := os.Stat(fileName); !os.IsNotExist(err) {
			f, err := os.Open(fileName)
			if err != nil {
				panic(err)
			}
			return Eval(bufio.NewReader(f), env)
		}
	}
	panic(fmt.Sprintf("Dependency %s not found", importFileString))
}

func DictCallable(args []Object, env *Env) Object {
	dictObject := DictObject{[]*DictObjectEntry{nil}}
	if len(args)%2 != 0 {
		panic("Expecetd key-value pairs are arguments")
	}

	for i := 0; i < len(args); i += 2 {
		dictObject.Set(args[i], args[i+1], env)
	}

	return dictObject
}

func GetStr(obj Object, env *Env) StringObject {
	operatorFunc, ok := GetSlot(obj, "__str__")
	if !ok {
		panic(fmt.Sprintf("__str__ slot not found on %s.", obj))
	}
	operatorCallable := operatorFunc.(CallableObject).Callable
	return operatorCallable([]Object{obj}, env).(StringObject)
}

func EnvCallable(args []Object, env *Env) Object {
	envDictObject := DictCallable([]Object{}, env)

	for symbol, object := range env.Objects {
		envDictObject.(DictObject).Set(StringObject{symbol}, object, env)
	}
	return envDictObject
}

func ListCallable(args []Object, env *Env) Object {
	return ListObject{args}
}
