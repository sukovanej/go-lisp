package main

import (
	. "github.com/sukovanej/go-lisp/interpreter"
	"testing"
)

func TestDict(t *testing.T) {
	env := GetMainEnv()
	dictObject := DictCallable([]Object{}, env)
	dictObject.GetSlots()["__set-item__"].(CallableObject).Callable([]Object{StringObject{"key"}, StringObject{"value"}, dictObject}, env)

	if dictObject.(DictObject).Len() != 1 {
		t.Errorf("There must be one item in the dict.")
	}
}
