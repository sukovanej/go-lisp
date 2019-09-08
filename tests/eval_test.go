package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"path/filepath"
	"strings"
	"testing"
)

func compareObject(left Object, right Object) bool {
	switch left.(type) {
	case NumberObject:
		return left.(NumberObject).Integer == right.(NumberObject).Integer
	case BoolObject:
		return left.(BoolObject).Value == right.(BoolObject).Value
	case NilObject:
		return left.(NilObject) == right.(NilObject)
	}

	return false
}

var expectedObject, outputObject Object
var input *bufio.Reader

func TestEvalSimple(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("1"))
	expectedObject = NumberObject{1}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalPlus(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("(+ 1 2)"))
	expectedObject = NumberObject{3}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalPlus2(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("(+ (+ 1 2) 3)"))

	expectedObject = NumberObject{6}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalLambda(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("((fn (x) (+ x 1)) 1)"))

	expectedObject = NumberObject{2}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalDefVariable(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("(set var 1)"))

	expectedObject = NumberObject{1}
	outputObject = Eval(input, env)
	outputEnvObject, _ := env.GetEnvSymbol("var")
	if !compareObject(outputObject, outputEnvObject) || !compareObject(outputEnvObject, expectedObject) {
		t.Errorf("%v != %v != %v.", expectedObject, outputObject, outputEnvObject)
	}
}

func TestEvalComparable(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("(== 1 2)"))

	expectedObject = BoolObject{false}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalIf(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("(if (== 1 1) 2 3)"))

	expectedObject = NumberObject{2}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}

	input = bufio.NewReader(strings.NewReader("(if (== 1 2) 2 3)"))

	expectedObject = NumberObject{3}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalAssert(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader("(!assert 1 1)"))

	expectedObject, _ = env.GetEnvSymbol("#nil")
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestEvalFactorial(t *testing.T) {
	env := GetMainEnv()
	input = bufio.NewReader(strings.NewReader(`(set fact (fn (x) (
		if (== x 1) 
			1 
			(* x (fact (- x 1))))))
		(fact 4)`))

	expectedObject = NumberObject{24}
	outputObject = Eval(input, env)
	if !compareObject(outputObject, expectedObject) {
		t.Errorf("%v != %v.", expectedObject, outputObject)
	}
}

func TestCustomTests(t *testing.T) {
	env := GetMainEnv()
	matches, err := filepath.Glob("./*.gisp")

	if err != nil {
		panic(err)
	}

	for _, file := range matches {
		fileEnv := &Env{map[string]Object{}, env}
		EvalFile(file, fileEnv)
	}
}
