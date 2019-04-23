/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/modbus"
	"github.com/jinzhu/gorm"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type ModSpawner struct {
	ComChan chan map[string]interface{}
}

type message struct {
	queue string
}

func main() {

	i := []int{1, 2, 3, 4}
	var tee ModSpawner
	tee.ComChan = make(chan map[string]interface{}, 2)
	tee.ComChan
	for a, v := range i {
		go tTestFuncest(tee)
	}

}

func TestFunc(tee ModSpawner) {

	fmt.Println(<-tee.ComChan)

}

func MultiModStart(db *gorm.DB) ([]*modbus.RTUClientHandler, []error) {

	handler := modbus.NewRTUClientHandler(serialParam.ComPort)
	handler.BaudRate = serialParam.BaudRate
	handler.DataBits = serialParam.DataBits
	handler.Parity = serialParam.Parity
	handler.StopBits = serialParam.StopBits
	handler.Timeout = time.Duration(serialParam.Timeout) * time.Second
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	err := handler.Connect()
	handlers = append(handlers, handler)
	retErrs = append(retErrs, err)
	return handlers, retErrs

}
