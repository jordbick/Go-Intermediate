package main

import (
	"fmt"
	"reflect"
)

func main() {
	// can assign functions to variables
	// anonymous function definition
	anon := func(msg string) bool {
		fmt.Println(msg)
		return true
	}

	res := anon("Hello anonymous function!")
	fmt.Println(res)
	fmt.Println("anon is of type:", reflect.TypeOf(anon))
}
