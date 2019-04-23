package main

import (
	"fmt"
	"sync"
	"time"
)

func saygo(s string, wg *sync.WaitGroup) {
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println(s)
	wg.Done()
}

func say(s string) {
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go saygo("world", &wg)
	say("hello")
	wg.Wait()

}
