package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	go send(ch)
	fmt.Println("\nblocking send goroutine...")
	// length of unbuffered channel will always be 0
	fmt.Printf("channel length: %v\n", len(ch))
	go receive(ch)
	time.Sleep(time.Second * 2)
}

// Takes a channel of string values
func send(ch chan string) {
	// Sending the string "message to the channel"
	ch <- "message"
}

// Delay to block the send goroutine because there is not a sender and receiver ready
// Following the 1 second the channel will then be unblocked
func receive(ch chan string) {
	time.Sleep(time.Second * 1)
	fmt.Println("send go routine unblocked")
	fmt.Printf("channel length: %v\n", len(ch))
	// print out value received
	fmt.Printf("received: %v\n", <-ch)
}
