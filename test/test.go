package main

import "fmt"
import "time"

func main() {
	c := make(chan int)

	go subrot(c)
	time.Sleep(1 * time.Second)
	select {
	case c <- 1:
	default:
	}
	fmt.Println("Hello")
	time.Sleep(2 * time.Second)
	// c <- 99
	time.Sleep(1 * time.Second)
}
func subrot(c chan int) {
	fmt.Println("In Subrot")
	// time.Sleep(5 * time.	)
	msg := <-c
	fmt.Println(msg)
}
