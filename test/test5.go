package main

import (
	"bytes"
	"fmt"
	"time"
)

func main() {
	var buffer bytes.Buffer
	a := "hello"
	b := "world"
	for i := 0; i < 10; i++ {

		fmt.Println(time.Now())
		buffer.WriteString(a)
		buffer.WriteString(b)

		e := buffer.String()
		fmt.Println(time.Now())
		d := a + b
		fmt.Println(time.Now())
		fmt.Println(d, e)
	}
}
