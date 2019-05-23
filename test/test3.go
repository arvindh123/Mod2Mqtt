package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

type MultiModInterface interface {
	read(id byte, wg2 *sync.WaitGroup, goroName string)
}

type MyRtuHandler struct {
	*modbus.RTUClientHandler
}

type MyTcpHandler struct {
	*modbus.TCPClientHandler
}

func main() {
	var wg2 sync.WaitGroup

	// handler = modbus.NewTCPClientHandler("localhost:502")
	// handler := &MyTcpHandler{&modbus.TCPClientHandler()} // Create empty varaible with modbus TCPClientHandler

	var handlers []*MyTcpHandler
	for i := 0; i < 20; i++ {
		handler := &MyTcpHandler{modbus.NewTCPClientHandler("localhost:502")}
		handler.Timeout = 5 * time.Second
		handlers = append(handlers, handler)
	}

	// handler := modbus.NewRTUClientHandler("\\\\.\\COM11")
	// handler := &MyRtuHandler{modbus.NewRTUClientHandler("\\\\.\\COM11")}
	// handler.BaudRate = 9600
	// handler.DataBits = 8
	// handler.Parity = "N"
	// handler.StopBits = 1

	// var handlers []*MyRtuHandler
	// for i := 0; i < 20; i++ {
	// 	handler := &MyRtuHandler{modbus.NewRTUClientHandler("\\\\.\\COM11")}
	// 	handler.BaudRate = 9600
	// 	handler.DataBits = 8
	// 	handler.Parity = "N"
	// 	handler.StopBits = 1
	// 	handler.Timeout = 5 * time.Second
	// 	handlers = append(handlers, handler)

	// }

	for i, handler := range handlers {
		fmt.Printf("%+v\n", handler)
		err := handler.Connect()
		_ = err

		// go handler.read(2, &wg2, strconv.Itoa(i))
		go ReadData(handler, 2, &wg2, strconv.Itoa(i))
		wg2.Add(1)

	}

	// for i := 0; i < 20; i++ {
	// 	go read(handler, 2, &wg2, strconv.Itoa(i))
	// 	wg2.Add(1)
	// }
	// go read(handler, 2, &wg2, "Second")
	// go TCPread(handler, 2, &wg2, "Second")
	// wg2.Add(1)
	// go read(handler, 1, &wg2, "Frist")
	// go TCPread(handler, 2, &wg2, "Frist")
	// wg2.Add(1)

	wg2.Wait()

}

func ReadData(inter MultiModInterface, id byte, wg2 *sync.WaitGroup, goroName string) {
	inter.read(id, wg2, goroName)
}

func (handler *MyRtuHandler) read(id byte, wg2 *sync.WaitGroup, goroName string) {
	defer handler.Close()
	handler.SlaveId = id
	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(13, 4)
	fmt.Printf("Goro - %s - Ambient Temperture %d\n", goroName, results)
	fmt.Println(err)
	results, err = client.ReadInputRegisters(0, 1)
	fmt.Printf("Goro - %s - Process Value %d\n", goroName, results)
	fmt.Println(err)
	wg2.Done()
}

func (handler *MyTcpHandler) read(id byte, wg2 *sync.WaitGroup, goroName string) {
	defer handler.Close()
	handler.SlaveId = id
	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(13, 4)
	fmt.Printf("Goro - %s - Ambient Temperture %d\n", goroName, results)
	fmt.Println(err)
	results, err = client.ReadInputRegisters(0, 1)
	fmt.Printf("Goro - %s - Process Value %d\n", goroName, results)
	fmt.Println(err)
	wg2.Done()
}
