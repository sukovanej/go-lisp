package interpreter

type FormObject struct {
	Callable func([]SyntaxValue, *Env) Object
}

func (o FormObject) GetSlots() map[string]Object {
	return map[string]Object{}
}

func createLambdaFunction(declared_args []SyntaxValue, body SyntaxValue, env *Env) CallableObject {
	return CallableObject{func(args []Object, _ *Env) Object {
		internal_env := &Env{map[string]Object{}, env}
		for i, declared_arg := range declared_args {
			switch declared_arg.(type) {
			case Token:
				internal_env.Objects[declared_arg.(Token).Symbol] = args[i]
			default:
				panic("not defined behaviour")
			}
		}
		return EvalSyntax(body, internal_env)
	}}
}

func CreateLambdaForm(args []SyntaxValue, env *Env) Object {
	if len(args) != 2 {
		panic("Wrong number of arguments")
	}

	return createLambdaFunction(args[0].([]SyntaxValue), args[1], env)
}
