package main

import (
	. "github.com/sukovanej/go-lisp/interpreter"
	"testing"
)

func TestToken(t *testing.T) {
	env := GetMainEnv()
	dictObject := DictCallable([]Object{}, env)
	dictObject.GetSlots()["__set-item__"].(CallableObject).Callable([]Object{dictObject, StringObject{"key"}, StringObject{"value"}}, env)

	if dictObject.(DictObject).Len() != 1 {
		t.Errorf("There must be one item in the dict.")
	}
}
