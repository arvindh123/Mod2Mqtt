package mod

import (
	"strconv"

	"github.com/goburrow/modbus"
)

func getRtuHandlerFromRtuHandlers(RtuHandlers []*modbus.RTUClientHandler, SerPort string) (*modbus.RTUClientHandler, modbus.Client) {
	for _, handler := range RtuHandlers {
		if SerPort == handler.Address {
			return handler, modbus.NewClient(handler)
		}
	}
	return nil, nil
}

func getAsciiHandlerFromAsciiHandlers(AsciiHandlers []*modbus.ASCIIClientHandler, SerPort string) (*modbus.ASCIIClientHandler, modbus.Client) {
	for _, handler := range AsciiHandlers {
		if SerPort == handler.Address {
			return handler, modbus.NewClient(handler)
		}
	}
	return nil, nil
}

func getTcpHandlerFromTcpHandlers(TcpHandlers []*modbus.TCPClientHandler, IpAdd string, Port int) (*modbus.TCPClientHandler, modbus.Client) {
	for _, handler := range TcpHandlers {
		// fmt.Println("tcpHandler.Adress - ", handler.Address)
		// fmt.Println("IpAdd+ \":\" + strconv.Itoa(Port)", IpAdd+":"+strconv.Itoa(Port))
		if IpAdd+":"+strconv.Itoa(Port) == handler.Address {
			return handler, modbus.NewClient(handler)
		}
	}
	return nil, nil
}
