package main

import "fmt"

func main() {
	// Invoking function directly
	func() {
		fmt.Println("Hello anonymous function")
	}()
}
