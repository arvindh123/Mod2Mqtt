package main

import (
	"context"
	"fmt"
	"time"
)

func Ctxfunc(ctx context.Context, cont int) {
	fmt.Println("Context", cont, "Started")
	for {

		select {
		case <-ctx.Done():
			fmt.Println("Context - ", cont, "Done")

			return
		default:
		}
	}
}

func main() {
	ctx := context.Background()
	//Derive a context with cancel
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)
	_ = ctxWithCancel

	go Ctxfunc(ctxWithCancel, 1)
	go Ctxfunc(ctxWithCancel, 2)
	go Ctxfunc(ctxWithCancel, 3)
	go Ctxfunc(ctxWithCancel, 4)

	time.Sleep(5000 * time.Microsecond)

	cancelFunction()
	time.Sleep(1 * time.Microsecond)
}
