package main

import (
	"fmt"
	"time"
	"reflect"
)

func main() {
	ab := make(chan int)
	cd := make(chan int)

	go test(ab, cd)

	time.Sleep(5 * time.Second)

	ab <- 5
	go test(ab, cd)
	cd <- 6
	time.Sleep(5 * time.Second)

}
func test(ab chan int, cd chan int) {
	select {
	case cmd := <-ab:
		fmt.Println("hello", cmd)
	case cmd := <-cd:
		fmt.Println("world", cmd)
		if cmd 
	}
}


func NewNonBlockingChan() *NonBlockingChan {

	send := make(chan interface{})
	recv := make(chan interface{})
	n    := &NonBlockingChan{Send: send,Recv: recv,}
	values := list.New()
	sendClosed := false

	
	go func() {
		const (
			incomingCase = iota
			outgoingCase
			numCases
		)

		for !sendClosed || items.Len() > 0 {
			cases := make([]reflect.SelectCase, numCases)
			cases[incomingCase].Dir = reflect.SelectRecv
			cases[outgoingCase].Dir = reflect.SelectSend

			if !sendClosed {
    			cases[incomingCase].Chan = reflect.ValueOf(send)
			}			

			if items.Len() > 0 {
				cases[outgoingCase].Chan = reflect.ValueOf(recv)
				cases[outgoingCase].Send = reflect.ValueOf(values.Front().Value)
			}

			i, val, ok := reflect.Select(cases)

			switch i {
			case incomingCase:
				if ok {
					values.PushBack(val.Interface())
				} else {
					sendClosed = true
				}
			case outgoingCase:
				values.Remove(values.Front())
			}
			close(recv)
		}
	}()
	return n
}