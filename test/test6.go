package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	var pay []map[string]interface{}
	payload1 := make(map[string]interface{})
	payload2 := make(map[string]interface{})
	var wg sync.WaitGroup

	var mqClient mqtt.Client
	var mqErr error
	_ = mqErr
	asd := "hello"
	payload1[asd] = 123

	payload1["ssid1"] = 123

	payload1["ssid2"] = 123

	pay = append(pay, payload1)
	payload2["ssid3"] = 123

	payload2["ssid4"] = 123
	pay = append(pay, payload2)
	mqClient, mqErr = MqConnect()

	if mqErr == nil {
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go PublishCon(mqClient, fmt.Sprintf("Goro-%d",i), pay, wg)
		}

	}

	wg.Wait()

}

func MqConnect() (mqtt.Client, error) {
	var er error
	clientid := fmt.Sprintf("TestClient_%d", rand.Intn(1000))
	server := "tcp://127.0.0.1:1883"
	username := "admin"
	password := "public"
	keepalive := 60
	opts := CreateMqClientOptions(clientid, server, username, password, keepalive)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c, er
}

func CreateMqClientOptions(clientid string, server string, username string, password string, keepalive int) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientid)
	opts.SetKeepAlive(2 * time.Second)
	return opts

}

func PublishCon(mqClient mqtt.Client, name string, pay []map[string]interface{}, wg sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		final, err := json.Marshal(pay)

		fmt.Println(final, err)
		fmt.Println(mqClient.Publish("hello/world/"+name, 0, false, final).Error())

		// mqClient.Publish("hello/world", 0, false, "payload-123-payload")
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}
