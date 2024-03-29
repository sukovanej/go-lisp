package interpreter

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func OperatorFunc(operatorName string) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		operatorFunc, ok := GetSlot(args[0], "__"+operatorName+"__")
		if !ok {
			return NewErrorWithoutToken(fmt.Sprintf("Slot %s not found on %v.", operatorName, args[0]))
		}
		operatorCallable := operatorFunc.(CallableObject).Callable
		result := args[0]
		for _, obj := range args[1:len(args)] {
			result = operatorCallable([]Object{result, obj}, env)
		}
		return result
	}
}

func Equal(left Object, right Object, env *Env) Object {
	operatorFunc, ok := GetSlot(left, "__==__")
	if !ok {
		return NewErrorWithoutToken(fmt.Sprintf("Slot __==__ not found on %v.", left))
	}
	operatorCallable := operatorFunc.(CallableObject).Callable
	return operatorCallable([]Object{left, right}, env)
}

func SetForm(args []Object, env *Env) Object {
	if len(args) != 2 {
		return NewErrorWithSyntaxValue(args[0].(SyntaxObject).Value, fmt.Sprintf("Form set expects 2 arguments, %d given.", len(args)))
	}

	obj := EvalSyntax(args[1].(SyntaxObject).Value, env)
	env.SetSymbol(args[0].(SyntaxObject).Value.(SymbolValue).Value.Symbol, obj)
	return obj
}

func AssertCallable(args []Object, env *Env) Object {
	if len(args) != 2 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable assert expects 2 arguments, %d given.", len(args)))
	}

	eqFunc, _ := env.GetEnvSymbol("==")
	equalConditionObject := eqFunc.(CallableObject).Callable(args, env)

	if IsErrorObject(equalConditionObject) {
		return equalConditionObject
	}

	if !equalConditionObject.(BoolObject).Value {
		stringLeft := GetStr(args[0], env)
		if IsErrorObject(stringLeft) {
			return stringLeft
		}
		stringRight := GetStr(args[1], env)
		if IsErrorObject(stringRight) {
			return stringRight
		}
		return NewErrorWithoutToken(fmt.Sprintf("Assertion error. %s != %s", stringLeft, stringRight))
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

func CreateStructForm(slots []Object, env *Env) Object {
	slotObjects := map[string]Object{}
	nilObject, _ := env.GetEnvSymbol("#nil")

	for _, slot := range slots {
		slotObjects[slot.(SyntaxObject).Value.(SymbolValue).Value.Symbol] = nilObject
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

func DefStructForm(declared_args []Object, env *Env) Object {
	constructor := CallableObject{func(args []Object, _ *Env) Object {
		structObject := CreateStructForm(declared_args[1:], env)
		structSlots := structObject.GetSlots()
		for i, _ := range declared_args[1:] {
			structSlots[declared_args[i+1].(SyntaxObject).Value.(SymbolValue).Value.Symbol] = args[i]
		}
		return structObject
	}}

	env.SetSymbol(declared_args[0].(SyntaxObject).Value.(SymbolValue).Value.Symbol, constructor)
	return constructor
}

func GetAttrForm(args []Object, env *Env) Object {
	if len(args) != 2 {
		panic("Unexpected number of arguments")
	}
	obj := EvalSyntax(args[0].(SyntaxObject).Value, env)
	slot, ok := GetSlot(obj, args[1].(SyntaxObject).Value.(SymbolValue).Value.Symbol)
	if !ok {
		panic(fmt.Sprintf("Slot %s not found", args[1].(SyntaxObject).Value.(SymbolValue).Value.Symbol))
	}
	return slot
}

func SetAttrForm(args []Object, env *Env) Object {
	if len(args) != 3 {
		panic("Unexpected number of arguments")
	}

	slots := EvalSyntax(args[0].(SyntaxObject).Value, env).GetSlots()
	_, ok := slots[args[1].(SyntaxObject).Value.(SymbolValue).Value.Symbol]
	if !ok {
		panic(fmt.Sprintf("Struct doesn't have slot '%s'.", args[1].(SyntaxObject).Value.(SymbolValue).Value.Symbol))
	}

	object := EvalSyntax(args[2].(SyntaxObject).Value, env)
	slots[args[1].(SyntaxObject).Value.(SymbolValue).Value.Symbol] = object
	return object
}

func NotEqualOperator(args []Object, env *Env) Object {
	result := OperatorFunc("==")(args, env).(BoolObject)
	return BoolObject{!result.Value}
}

func SlotCallable(slotName string, numberOfArguments int) func([]Object, *Env) Object {
	return func(args []Object, env *Env) Object {
		if len(args) != numberOfArguments {
			return NewErrorWithoutToken(fmt.Sprintf("Callable expects %d arguments, %d was given.", numberOfArguments, len(args)))
		}
		operatorFunc, ok := GetSlot(args[len(args)-1], slotName)
		if !ok {
			return NewErrorWithoutToken(fmt.Sprintf("%s slot not found.", slotName))
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
			return Eval(bufio.NewReader(f), env, &BufferMetaInformation{0, 0, fileName})
		}
	}
	panic(fmt.Sprintf("Dependency %s not found", importFileString))
}

func DictCallable(args []Object, env *Env) Object {
	dictObject := DictObject{map[string]DictObjectEntry{}}
	if len(args)%2 != 0 {
		panic("Expecetd key-value pairs are arguments")
	}

	for i := 0; i < len(args); i += 2 {
		dictObject.Set(args[i], args[i+1], env)
	}

	return dictObject
}

func GetStr(obj Object, env *Env) Object {
	operatorFunc, ok := GetSlot(obj, "__str__")
	if !ok {
		return NewErrorWithoutToken(fmt.Sprintf("__str__ slot not found on %s.", obj))
	}
	operatorCallable := operatorFunc.(CallableObject).Callable
	return operatorCallable([]Object{obj}, env)
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

func AndForm(args []Object, env *Env) Object {
	if len(args) == 0 {
		panic("Wrong number of arguments")
	}

	for _, arg := range args {
		value := EvalSyntax(arg.(SyntaxObject).Value, env)
		if !value.(BoolObject).Value {
			f, _ := env.GetEnvSymbol("#f")
			return f
		}
	}

	t, _ := env.GetEnvSymbol("#t")
	return t
}

func OrForm(args []Object, env *Env) Object {
	if len(args) == 0 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable or expects at least 1 argument, 0 given."))
	}

	for _, arg := range args {
		value := EvalSyntax(arg.(SyntaxObject).Value, env)

		if IsErrorObject(value) {
			return value
		}

		if !IsBoolObject(value) {
			return NewErrorWithoutToken(fmt.Sprintf("Object must be bool."))
		}

		if value.(BoolObject).Value {
			t, _ := env.GetEnvSymbol("#t")
			return t
		}
	}

	f, _ := env.GetEnvSymbol("#f")
	return f
}

func GreaterOperator(args []Object, env *Env) Object {
	smaller := OperatorFunc("<")(args, env).(BoolObject).Value
	equal := OperatorFunc("==")(args, env).(BoolObject).Value
	return BoolObject{!smaller && !equal}
}

func NotOperator(args []Object, env *Env) Object {
	if len(args) != 1 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable not expects 1 argument, %d given.", len(args)))
	}

	if !IsBoolObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be bool."))
	}
	return BoolObject{!args[0].(BoolObject).Value}
}

func ShellCommandCallable(args []Object, env *Env) Object {
	if len(args) != 1 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable not expects 1 argument, %d given.", len(args)))
	}

	if !IsStringObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be string."))
	}

	cmd := exec.Command("/bin/sh", []string{"-c", args[0].(StringObject).String}...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return NewErrorWithoutToken(err.Error())
	}
	return StringObject{out.String()}
}

func ApplyCallable(args []Object, env *Env) Object {
	if len(args) != 2 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable not expects 2 argument, %d given.", len(args)))
	}

	if !IsCallableObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be callable."))
	}

	if !IsListObject(args[1]) {
		fmt.Println(args[1])
		return NewErrorWithoutToken(fmt.Sprintf("Object must be list."))
	}

	return args[0].(CallableObject).Callable(args[1].(ListObject).List, env)
}

func CommandLineArgumentsCallable(args []Object, env *Env) Object {
	listObjects := []Object{}
	for _, arg := range os.Args {
		listObjects = append(listObjects, StringObject{arg})
	}
	return ListObject{listObjects}
}

func PrognForm(args []Object, env *Env) Object {
	last, _ := env.GetEnvSymbol("#nil")
	for _, statement := range args {
		last := EvalSyntax(statement.(SyntaxObject).Value, env)

		if IsErrorObject(last) {
			return last
		}
	}
	return last
}

func EvalCallable(args []Object, env *Env) Object {
	if len(args) != 1 && len(args) != 2 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable eval expects 2 or 1 arguments, %d given.", len(args)))
	}

	if !IsSyntaxObject(args[0]) {
		return NewErrorWithoutToken(fmt.Sprintf("Object must be syntax object."))
	}

	return EvalSyntax(args[0].(SyntaxObject).Value, env)
}

func QuoteForm(args []Object, env *Env) Object {
	if len(args) != 1 {
		return NewErrorWithoutToken(fmt.Sprintf("Callable eval expects 1 argument, %d given.", len(args)))
	}

	return args[0]
}
