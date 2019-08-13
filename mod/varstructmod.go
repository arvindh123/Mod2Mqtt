package mod

import (
	"sync"

	"github.com/goburrow/modbus"
)

type RtuContainer struct {
	RTUClientHandler *modbus.RTUClientHandler
	RTUClient        modbus.Client
}

type TcpContainer struct {
	TCPClientHandler *modbus.TCPClientHandler
	TCPClient        modbus.Client
}

var processMu sync.Mutex
var AddFeatures map[string]interface{}
var RtuHandlers []*modbus.RTUClientHandler
var TcpHandlers []*modbus.TCPClientHandler
var wg sync.WaitGroup
var modStatus bool = false
