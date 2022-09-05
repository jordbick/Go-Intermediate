package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	fmt.Println("Running in main routine")

	go g1()
	time.Sleep(10 * time.Millisecond)

	fmt.Println("exited from main routine")

	// ANONYMOUS GO ROUTINES

	for i := 0; i < 5; i++ {
		// closure as go routine is accessing variable declared outside of routine
		// Each iteration of the loop is a new Goroutine closure - all the closures are using the same instance of the variable i
		// Each closure prints the value of i and by the time they get round to printing the value of i, the loop had actually already completed and the exit value for the loop is i = 5
		// So they all display the value of 5
		go func() {
			fmt.Println("i = ", i)
		}()
	}
	time.Sleep(10 * time.Millisecond)

	// To solve this problem, we can define the closure with an integer parameter and pass the closure an argument when its invoked
	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println("i = ", i)
		}(i)
	}
	time.Sleep(10 * time.Millisecond)

	// Can also assign the function to a variable and invoke the function variable
	for i := 0; i < 5; i++ {
		f := func(i int) {
			fmt.Println("i = ", i)
		}
		go f(i)
	}
	time.Sleep(10 * time.Millisecond)

	// When using closure we have to keep in mind that the closure can reference variables declared in the surrounding function
	// Each Goroutine is going to have acces to the same variable value i
	// So by refactoring our code, we are now making it so each Goroutine gets its own isolated value for i
	// They are not displayed in order because the routines are running concurrently

	// SYNC AND SYNC/ATOMIC

	// sync package includes some basic synchronisation primitives, such as Mutex, watgroups etc
	// sync.Mutex = Allows us to use a mutual exlcusion on a shared resource
	// sync.RWMutex = Reader/Writer mutual exclusion lock. Enables a lock to be held by any number of concurrent reader or a single writer
	// sync.WaitGroup = Wait for goroutines to complete
	// sync.Map = Similar to a regular map but is concurrency safe
	// sync.Pool = Concurrent pool used to manage the creation of expensive resources like DB or network connections
	// sync.Once = Run a one-time initialisation function, commonly used to call initialisation functions that multiple Gorutines depend on
	// sync.Cond = Implements a condition variable. Provides a vehicle for Goroutines to wait for, or announce, that some event has occurred

	// Atomic operations - More primitive that other sync methods
	// Lockless, usually implemented at the hardware level and sometimes used by other synchronisation methods
	// package contains low-level memory primitives for synchronisation
	// Used for low-level applications
	// LoadPointer
	// StorePointer
	// SwapPointer
	// CompareAndSwapPointer

	// WAITGROUP
	// Defined the wg variable
	var wg sync.WaitGroup
	// Use Add method to add a count for each Goroutine that gets spawned
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		// When we call the goroutine we pass the pointer to the WaitGroup variable
		// Need a pointer otherwise it will be a copy everytime and won't work
		go myfunc(&wg, i)
	}
	// Call to Wait to wait for the Goroutines to complete
	wg.Wait()
	fmt.Println("Each goroutine has to run to completion, thanks for waiting!")

	// ATOMIC OPERATIONS

	// Set GOMAXPROCS to use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	var counter uint64
	var wg1 sync.WaitGroup

	// For each iteration 1 added to WaitGroup and spawn a goroutine
	for i := 0; i < 1000; i++ {
		wg1.Add(1)
		go func() {
			// Each goroutine includes a defer call to the Done method and increments the counter 1000 times
			defer wg1.Done()
			for j := 0; j < 1000; j++ {
				// counter++
				// This is now incremented in a concurrency safe manner
				atomic.AddUint64(&counter, 1)
			}
		}()
	}
	// Wait for all the go routines to execute
	wg1.Wait()
	// Will have wrong value (sould be result of 1 million)
	// However the counter variable is being accessed by multiple Goroutines that run concurrently
	// The increment operation is not conducted in a concurrency safe manner right now
	fmt.Println("counter: ", counter)

	// MUTEXS
	var counter1 int
	var wg2 sync.WaitGroup

	// We are using shared variables so the concurrent routines can access the variable at differnt times
	// Need to implement a lock so it can't be written to by many different functions

	var mu sync.Mutex

	// variables assigned anonymous functions
	inc := func() {
		// Need to lock and then unlock the mutex, always have the unlock deferred right after the lock
		mu.Lock()
		defer mu.Unlock()
		counter1++
	}

	dec := func() {
		mu.Lock()
		defer mu.Unlock()
		counter1--
	}

	// For each iteration adds one to the wait group, sets up a goroutine (as anon function), defers Done and increments
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			inc()
		}()
	}

	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			dec()
		}()
	}

	wg2.Wait()
	fmt.Println("final value of counter", counter1)

	// RACE CONDITIONS
	// A condition whereby a program depends on the sequence or timing of its processes or threads, and that timing cannot be guaranteed
	// Can occur in concurrent programmes that access and modify a shared resource
	// Occurs when there is unsynchronised access to shared memory

	// DATA RACES
	// Conditions where mutliple threads are accessing shared memory concurrently
	// At least one of those threads is performing a write operation
	// Resulting in one or more of the threads being unsynchronised with the other threads
	// Can block with waitgroups, channels, return channels to push a result or use mutexs
	// Can detect race conditions using race detector "go run -race main.go"

	// GO CHANNELS
	// How goroutines communicate
	// Created using "make" keyword
	// Used to send a receive data of a specific type
	// When a channel is created using make, it allocates memory for the channel and returns a reference to it
	// Channel operator -> <-
	// ch := make(chan int) // channel of ints, can also set the buffer capacity of the channel. If not set it is an unbuffered channel
	// ch <- val // send int val to ch
	// val := <- ch // data received from the channel ch and written to the variable val
	// Allows synchronisation by blocking without having to explicity use locks or condition variables

	// Buffered channels have a buffer between the send and receiver goroutines - designates the number of elements that can be sent without the receiver being ready to receive
	// If unbuffered, the channel is synchronous

	ch := make(chan int)
	mult := func(x, y int) {
		res := x * y
		// res being sent to channel ch
		ch <- res
	}

	go mult(10, 10)
	// Value val is being received from the channel ch
	// when you have a receive operation also returns bool that signifies whether the value received from the channel was generated by a write (true) or by a close (false) operation on the channel
	val, ok := <-ch

	fmt.Printf("Type of value of ch: %T\n", ch) // channel of ints
	fmt.Printf("value of ch: %v\n", ch)         // memory address
	fmt.Printf("Resulting value from goroutine: %v \n", val)
	fmt.Printf("value of ok: %v \n", ok)

	// CLOSE AND RANGE METHODS ON CHANNELS
	// Using the range function, we can range over the values in a channel

	ch1 := make(chan string, 2)
	// goroutine which accepts a string channel. The sender
	go func(ch chan string) {
		// guarentees we close the channel when the anon function finishes
		defer close(ch)
		for i := 1; i <= 5; i++ {
			msg := "message" + strconv.Itoa(i)
			// write the message to the channel
			ch <- msg
			fmt.Printf("SEND goroutine: %v\n", msg)
		}
	}(ch1)

	// goroutine which is the receiver
	go func(ch chan string) {
		// for loop to range over values in the channel
		// range will automatically break out of the loop when the channel is closed
		for val := range ch {
			fmt.Printf("RECV go routine: %v\n", val)
		}
	}(ch1)

	fmt.Println("MAIN go routine: sleeping!")
	time.Sleep(time.Millisecond * 100)
	//
	val1, ok := <-ch1
	fmt.Printf("values returned from the read operation 'val, ok := <-ch1' are: %q, %v\n", val1, ok)
	fmt.Println("MAIN goroutine: done!")

	// CHANNEL BLOCKING
	// By default channels are blocking
	// Select case statements can be used as non-blocking sends and receives

	// PIPELINES
	// Specific type of concurrent program that consists of stages connected by channels
	// Each stage is a group of goroutines running the same function, isolating processing from stage to stage

	vals := []int{100, 50, 20, 90}
	// gen function is a stage that takes a slice of ints and copies them to an output channel
	in := gen(vals)

	// fan-out
	// Split up work across 2 channels. Call to square will make a goroutine
	fo1 := square(in)
	fo2 := square(in)

	// fan-in
	// merge stage takes a list of channels and copies the values to a single merged channel
	fi := merge(fo1, fo2)
	for res := range fi {
		fmt.Println(res)
	}

}

func g1() {
	fmt.Println("running in goroutine g1")
}

// Pass in a refernce to the WaitGroup type
func myfunc(wg *sync.WaitGroup, i int) {
	// Defer the call to wg.Done(), so when each goroutine completes it decrements the counter by 1 till it reaches 0
	// And the main block is no longer blocked and can exit
	// Can use sync/atomic to help with this
	defer wg.Done()
	fmt.Println("Finished executing iteration", i)
}

// take slice of ints, coppy to an out channel and then return the out channel to the calling routine
func gen(vals []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, val := range vals {
			out <- val
		}
		close(out)
	}()
	return out
}

// stage that takes a receive channel as input, and for each value in the channel it multiplies the value by itself
// then copies those values to an output channel named out
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for val := range in {
			out <- val * val
		}
		close(out)
	}()
	return out
}

// Takes a variable number of channels and then output the values on one channel
func merge(fo ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// anon function which takes a channel and copies the values on a channel to a new output channel
	fi := func(ch <-chan int) {
		for val := range ch {
			out <- val
		}
		wg.Done()
	}

	// Add to the waitgroup counter the number of channels
	wg.Add(len(fo))
	for _, ch := range fo {
		go fi(ch)
	}

	// once wg down to 0 close channel
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
