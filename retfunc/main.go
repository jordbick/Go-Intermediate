package main

import (
	"fmt"
	"reflect"
)

func main() {
	// assigned myfunc the return value of retfunc, because we are invoking the function. We are returned an anonymous function
	myfunc := retfunc()
	// Here we are invoking the return value (the anonymous function) so returns true
	fmt.Println(myfunc())
	fmt.Println("myfunc is of type", reflect.TypeOf(myfunc))
	fmt.Println("myfunc return value is of type", reflect.TypeOf(myfunc()))

}

// Declared function called retfunc that two return types -  a func() and a bool inside that function
func retfunc() func() bool {
	// returning an anon function
	// inside the anon function we are returning the bool value true
	return func() bool {
		return true
	}
}
