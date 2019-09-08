package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

func compareObject(left Object, right Object) bool {
	switch left.(type) {
	case NumberObject:
		return left.(NumberObject).Integer == right.(NumberObject).Integer
	}

	return false

	//leftSlots := left.GetSlots()
	//rightSlots := left.GetSlots()

	//for k, _ := range leftSlots {
	//	if !compareObject(leftSlots[k], rightSlots[k]) {
	//		return false
	//	}
	//}
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
