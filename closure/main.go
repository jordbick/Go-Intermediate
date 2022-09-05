package main

import (
	"fmt"
	"reflect"
)

// A closure is an anon function that can reference variables that are declared external to the anon function within a surrounding function
// The variables declared in the surrounding function can be shared between the closure and the surrounding function

func main() {
	// assigning x the return value of invoking inc() = anon function or closure
	x := inc()
	fmt.Println("x is of type:", reflect.TypeOf(x))

	for i := 0; i < 4; i++ {
		// inside the function x() we are incrementing the value j which starts out at 0
		fmt.Println("x = ", x())
	}

	// access to its own copy of inc()
	y := inc()
	fmt.Println("y = ", y())
}

// returns an anon function and an int
func inc() func() int {
	var j int
	// closure returning a variable declared outside the anon function
	return func() int {
		j++
		return j
	}
}
