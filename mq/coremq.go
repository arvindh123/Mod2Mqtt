package mq

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/arvindh123/Mod2Mqtt/ws"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqClient mqtt.Client
var mqErr error
var mqttStatus bool = false

func MqConnect(CidPart string, mqParams models.MqttDetails) (bool, error) {
	clientid := fmt.Sprintf("%s_%d", CidPart, rand.Intn(2147483647))
	server := fmt.Sprintf("tcp://"+mqParams.Ip+":%d", mqParams.Port)
	username := mqParams.Username
	password := mqParams.Password
	keepalive := 60
	opts := CreateMqClientOptions(clientid, server, username, password, keepalive)
	mqClient = mqtt.NewClient(opts)
	token := mqClient.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		go ws.WsStatusPub(fmt.Sprintf("Error in Mqtt Connection - %v", err))
		mqttStatus = false
		return mqttStatus, err
	}
	mqttStatus = true

	return mqttStatus, nil
}

func connLostHandler(c mqtt.Client, err error) {
	go ws.WsStatusPub(fmt.Sprintf("%v	MQTT Server Connection lost, reason: %v\n", time.Now(), err))
	mqttStatus = false
}

func onConnectHandler(c mqtt.Client) {
	go ws.WsStatusPub(fmt.Sprintf("%v	MQTT Server Connected", time.Now()))
	mqttStatus = true
}

func CreateMqClientOptions(clientid string, server string, username string, password string, keepalive int) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientid)
	opts.SetConnectTimeout(time.Duration(20) * time.Second)
	opts.SetKeepAlive(time.Duration(60) * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connLostHandler)
	opts.SetOnConnectHandler(onConnectHandler)
	return opts
}

func MqttProcStop() bool {
	if mqClient != nil {
		mqClient.Disconnect(250)
		mqttStatus = false
	}
	return mqttStatus
}
func MqttPublish2(Device models.DeviceDetails, Payload chan map[string]interface{}, MqWg *sync.WaitGroup) {
	if mqttStatus == true {
		finalPay := make(map[string]interface{})
		// fmt.Println("Starting <- Pays")
		for vs := range Payload {
			// fmt.Println("for each payload - ", vs)
			for k, v := range vs {
				finalPay[k] = v
			}
		}
		if len(finalPay) > 0 {
			final, err := json.Marshal(finalPay)
			// fmt.Println("Pusblished data %v", finalPay)

			if err == nil {
				mqClient.Publish(Device.PubTopic, 0, false, final)
				go ws.WsClientPub(ws.MqLastSent{fmt.Sprintf("Topic - %s, Payload- %v", Device.PubTopic, finalPay)})
			} else {
				go ws.WsClientPub(ws.MqLastSent{fmt.Sprintf("Topic - %s, Before Marshal Payload- %v, Error in Marshal -> Err- %v", Device.Name, finalPay, err)})
			}
		}
	} else if mqttStatus == false {
		go ws.WsClientPub(ws.MqLastSent{fmt.Sprintf("MQTT Server Connection Lost")})
	}
	MqWg.Done()
	return
}
