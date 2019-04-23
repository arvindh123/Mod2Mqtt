package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	var wg2 sync.WaitGroup

	handler := modbus.NewRTUClientHandler("\\\\.\\COM11")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1

	handler.Timeout = 5 * time.Second
	fmt.Printf("%+v\n", handler)
	err := handler.Connect()
	_ = err
	defer handler.Close()

	go read(handler, 2, &wg2, "Second")
	wg2.Add(1)
	// go read(handler, 1, &wg2, "Frist")
	// wg2.Add(1)

	wg2.Wait()

}

func read(handler *modbus.RTUClientHandler, id byte, wg2 *sync.WaitGroup, goroName string) {
	handler.SlaveId = id
	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(13, 4)
	fmt.Printf("Goro - %s - Ambient Temperture %d", goroName, results)
	fmt.Println(err)
	results, err = client.ReadInputRegisters(0, 1)
	fmt.Printf("Goro - %s - Process Value %d ", goroName, results)
	fmt.Println(err)
	wg2.Done()
}
