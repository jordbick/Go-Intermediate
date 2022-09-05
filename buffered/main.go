package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	// capacity of channel
	ch := make(chan string, 3)
	go send(ch)
	fmt.Println("\npress enter key to continue...")
	// Scanln to pause program until enter pressed
	fmt.Scanln()
	go receive(ch)
	time.Sleep(time.Second * 1)
}

// passes messsages over channel and measures length of channel
func send(ch chan string) {
	fmt.Printf("channel length before sending anything: %v\n", len(ch))
	for i := 1; i < 4; i++ {
		msg := "message" + strconv.Itoa(i)
		fmt.Printf("sending: %v\n", msg)
		ch <- msg
		fmt.Printf("channel length: %v\n", len(ch))
	}

}

func receive(ch chan string) {
	for i := 0; i < 3; i++ {

		// print out value received
		fmt.Printf("received: %v\n", <-ch)
		fmt.Printf("channel length: %v\n", len(ch))
	}
}
