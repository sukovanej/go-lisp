package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func EvalFile(file string, env *Env) Object {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	meta := &BufferMetaInformation{0, 0, file}
	return Eval(bufio.NewReader(f), env, meta)
}

func SetupMainEnv() *Env {
	env := GetMainEnv()
	systemPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	pwdPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	env.SetSymbol("__import_path__", ListObject{[]Object{
		StringObject{systemPath + "/std"},
		StringObject{pwdPath},
	}})
	ImportCallable([]Object{StringObject{"builtin"}}, env)
	return env
}

func Eval(reader *bufio.Reader, env *Env, meta *BufferMetaInformation) Object {
	syntaxTree := GetSyntax(reader, meta)
	var lastResult Object

	for syntaxTree != nil {
		lastResult = EvalSyntax(syntaxTree, env)

		if IsErrorObject(lastResult) {
			return lastResult
		}

		syntaxTree = GetSyntax(reader, meta)
	}

	return lastResult
}

func EvalSyntax(value SyntaxValue, env *Env) Object {
	switch value.GetType() {
	case SYNTAX_SYMBOL:
		return evalSymbol(value.(SymbolValue).Value, env)
	case SYNTAX_LIST:
		return evalFunction(value.(ListValue).Value, env)
	case SYNTAX_LIST_LITERAL:
		return evalListLiteral(value.(ListLiteralValue), env)
	case SYNTAX_DICT_LITERAL:
		return evalDictLiteral(value.(DictLiteralValue), env)
	}

	panic("Unexpected syntax token.")
}

func evalListLiteral(args ListLiteralValue, env *Env) Object {
	argObjects := []Object{}
	for _, value := range args.Value {
		argObjects = append(argObjects, EvalSyntax(value, env))
	}
	return ListCallable(argObjects, env)
}

func evalDictLiteral(args DictLiteralValue, env *Env) Object {
	argObjects := []Object{}
	for _, value := range args.Value {
		argObjects = append(argObjects, EvalSyntax(value, env))
	}
	return DictCallable(argObjects, env)
}

func evalSymbol(token Token, env *Env) Object {
	if token.Type == SYMBOL_STRING {
		return StringObject{token.Symbol}
	} else if num, err := strconv.Atoi(token.Symbol); err == nil {
		return NumberObject{num}
	} else {
		obj, ok := env.GetEnvSymbol(token.Symbol)
		if !ok {
			return NewError(token, fmt.Sprintf("Symbol %s not found", token.Symbol))
		}
		return obj
	}
}

func evalFunction(list []SyntaxValue, env *Env) Object {
	function := EvalSyntax(list[0], env)

	if IsErrorObject(function) {
		return function
	}

	switch function.(type) {
	case CallableObject:
		args := []Object{}

		for _, arg := range list[1:len(list)] {
			value := EvalSyntax(arg, env)
			if IsErrorObject(value) {
				return value
			}
			args = append(args, value)
		}

		return function.(CallableObject).Callable(args, env)
	case FormObject:
		args := list[1:len(list)]
		return function.(FormObject).Callable(args, env)
	default:
		return NewErrorWithSyntaxValue(list[0], "First item must be callable")
	}
}

func GetMainEnv() *Env {
	return &Env{
		map[string]Object{
			"+":         CallableObject{OperatorFunc("+")},
			"-":         CallableObject{OperatorFunc("-")},
			"*":         CallableObject{OperatorFunc("*")},
			"/":         CallableObject{OperatorFunc("/")},
			"^":         CallableObject{OperatorFunc("^")},
			"fn":        FormObject{CreateLambdaForm},
			"set":       FormObject{SetForm},
			"#t":        BoolObject{true},
			"#f":        BoolObject{false},
			"==":        CallableObject{OperatorFunc("==")},
			"<":         CallableObject{SlotCallable("__<__", 2)},
			"!=":        CallableObject{NotEqualOperator},
			">":         CallableObject{GreaterOperator},
			"#nil":      NilObject{},
			"if":        FormObject{IfForm},
			"!assert":   CallableObject{AssertCallable},
			"print":     CallableObject{PrintCallable},
			"println":   CallableObject{PrintlnCallable},
			"->":        FormObject{GetAttrForm},
			"set->":     FormObject{SetAttrForm},
			"struct":    FormObject{CreateStructForm},
			"defn":      FormObject{DefLambdaForm},
			"str":       CallableObject{SlotCallable("__str__", 1)},
			"defstruct": FormObject{DefStructForm},
			"dict":      CallableObject{DictCallable},
			"item":      CallableObject{SlotCallable("__item__", 2)},
			"set-item":  CallableObject{SlotCallable("__set-item__", 3)},
			"env":       CallableObject{EnvCallable},
			"import":    CallableObject{ImportCallable},
			"list":      CallableObject{ListCallable},
			"len":       CallableObject{SlotCallable("__len__", 1)},
			"slice":     CallableObject{SlotCallable("__slice__", 3)},
			"append":    CallableObject{SlotCallable("__append__", 2)},
			"and":       FormObject{AndForm},
			"or":        FormObject{OrForm},
			"cond":      FormObject{CondForm},
			"not":       CallableObject{NotOperator},
		},
		nil,
	}
}
