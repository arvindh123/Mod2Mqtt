package main

import (
	"fmt"
	"time"

	"github.com/goburrow/modbus"
)

func main() {

	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(13, 1)
	fmt.Printf("Ambient Temperture %d", results)
	fmt.Println(err)
	results, err = client.ReadHoldingRegisters(31, 1)
	fmt.Printf("Process Value %d ", results)
	fmt.Println(err)

}
