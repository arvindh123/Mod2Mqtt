package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type StatusLog struct {
	StatusLog string `json:"statuslog"`
}

type ModLastAquired struct {
	ModLastAquired string `json:"modlastaquired"`
}
type MqLastSent struct {
	MqLastSent string `json:"mqlastsent"`
}

var WsClients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func WsClientPub(msg interface{}) {
	mu.Lock()
	defer mu.Unlock()
	wsJson, err := json.Marshal(msg)
	if err == nil {
		for client := range WsClients {
			// fmt.Println("--------------", msg)
			err := client.WriteMessage(1, wsJson)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(WsClients, client)
			}
		}
	}
}

func WsStatusPub(msg string) {
	mu.Lock()
	defer mu.Unlock()
	wsJson, err := json.Marshal(StatusLog{msg})
	if err == nil {
		for client := range WsClients {
			err := client.WriteMessage(1, wsJson)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(WsClients, client)
			}
		}
	}
}
