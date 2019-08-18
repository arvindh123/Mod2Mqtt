package mod

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/arvindh123/Mod2Mqtt/mq"
	"github.com/arvindh123/Mod2Mqtt/ws"
	"github.com/goburrow/modbus"
	// "github.com/goburrow/modbus"
)

func ModProcStart(Inters models.AllStructParams, ctx context.Context, addFeatures map[string]interface{}) bool {
	if modStatus == false {
		processMu.Lock()
		defer processMu.Unlock()
		// ctx, cancel = context.WithCancel(context.Background()
		AddFeatures = addFeatures

		runCount := 0
		for _, Inter := range Inters.AllStructParams {
			if Inter.Inter.Type == 1 {
				handler := modbus.NewRTUClientHandler(Inter.Inter.ComPort)
				handler.BaudRate = Inter.Inter.BaudRate
				handler.DataBits = Inter.Inter.DataBits
				handler.Parity = Inter.Inter.Parity
				handler.StopBits = Inter.Inter.StopBits
				if Inter.Inter.Timeout > 0 {
					handler.Timeout = 1000 * time.Millisecond
				} else {
					handler.Timeout = time.Duration(Inter.Inter.Timeout) * time.Millisecond
				}
				// handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
				err := handler.Connect()

				if err == nil {
					var Rtu RtuContainer
					Rtu.RTUClientHandler, Rtu.RTUClient = handler, modbus.NewClient(handler)
					wg.Add(1)
					runCount = runCount + 1
					RtuHandlers = append(RtuHandlers, handler)
					go Rtu.ModReadWrite(ctx, Inter, &wg)
				} else {
					handler.Close()
					go ws.WsStatusPub(fmt.Sprintf("%v	Error in Modstart - %s,  Serial Port - %s ", time.Now(), err.Error(), handler.Address))
				}

			} else if Inter.Inter.Type == 3 {
				var b bytes.Buffer
				b.WriteString(Inter.Inter.IpAdd)
				b.WriteString(":")
				b.WriteString(strconv.Itoa(Inter.Inter.Port))

				handler := modbus.NewTCPClientHandler(b.String())
				if Inter.Inter.Timeout > 0 {
					handler.Timeout = 1000 * time.Millisecond
				} else {
					handler.Timeout = time.Duration(Inter.Inter.Timeout) * time.Millisecond
				}
				err := handler.Connect()
				if err == nil {
					var Tcp TcpContainer
					Tcp.TCPClientHandler, Tcp.TCPClient = handler, modbus.NewClient(handler)
					wg.Add(1)
					runCount = runCount + 1
					TcpHandlers = append(TcpHandlers, handler)
					go Tcp.ModReadWrite(ctx, Inter, &wg)
				} else {
					handler.Close()
					go ws.WsStatusPub(fmt.Sprintf("%v	Error in Modstart - %s,  TCP - %s ", time.Now(), err.Error(), handler.Address))
				}
				// handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
			}

		}
		if runCount > 0 {
			modStatus = true
		} else {
			modStatus = false
		}
	} else {
		go ws.WsStatusPub(fmt.Sprintf("%v	Already Modbus Process running", time.Now()))
	}

	return modStatus
}

func ModProcStop(cancel context.CancelFunc) bool {
	processMu.Lock()
	defer processMu.Unlock()

	// go ws.WsStatusPub(fmt.Sprintf("%v	No of Gorotuines - %d", time.Now(), runtime.NumGoroutine()))
	go ws.WsStatusPub(fmt.Sprintf("%v	Topic Spawner and Modbus Ports closing", time.Now()))

	if cancel != nil {
		cancel()
	}

	go ws.WsStatusPub(fmt.Sprintf("%v	Waiting for Modbus WG", time.Now()))
	wg.Wait()
	go ws.WsStatusPub(fmt.Sprintf("%v	Finished all Modbus WG", time.Now()))

	for _, handler := range RtuHandlers {
		handler.Close()
	}
	for _, handler := range TcpHandlers {
		handler.Close()
	}

	go ws.WsStatusPub(fmt.Sprintf("%v	All Modbus Ports Closed", time.Now()))
	modStatus = false
	// go ws.WsStatusPub( fmt.Sprintf("%v	No of Gorotuines After - %d", time.Now(), runtime.NumGoroutine()))
	return modStatus
}

func (r *RtuContainer) ModReadWrite(ctx context.Context, inter models.InterWithDevices, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			// fmt.Println("Completed RTU ctx.Done()", inter.Inter.Name)
			wg.Done()
			return
		case <-time.After(time.Duration(inter.Inter.DaqRate) * time.Millisecond):
			for _, device := range inter.Devices {
				r.ModOperation(device)
			}
		}
	}

}

func (t *TcpContainer) ModReadWrite(ctx context.Context, inter models.InterWithDevices, wg *sync.WaitGroup) {
	runtime.GOMAXPROCS(50)
	var TcpDevWg sync.WaitGroup
	for _, device := range inter.Devices {
		TcpDevWg.Add(1)
		if inter.Inter.DaqRate > 0 {
			go t.ModTrigger(device, ctx, &TcpDevWg, inter.Inter.DaqRate)
		} else {
			go t.ModOperationSync(device, ctx, &TcpDevWg, inter.Inter.DaqRate)
		}

	}
	for {
		select {
		case <-ctx.Done():
			// fmt.Println("Completed ctx.Done()", inter.Inter.Name)
			TcpDevWg.Wait()
			wg.Done()
			return
		}
	}

}

func (r *RtuContainer) ModOperation(Device models.DeviceWithRegs) {
	Payload := make(chan map[string]interface{}, len(Device.Modregs))
	tpay := make(map[string]interface{})
	var ModWg sync.WaitGroup
	var MqWg sync.WaitGroup

	MqWg.Add(1)                                       // remove MqWait once archive to DB is done in MQTT side
	go mq.MqttPublish2(Device.Device, Payload, &MqWg) // change to go mq.MqttPublish2(Device.Device, Payload)" once archive to DB is done in MQTT side
	r.RTUClientHandler.SlaveId = Device.Device.MbId

	results := []byte{}
	err := errors.New("")
	if val, ok := AddFeatures["sts"]; ok {
		tpay["sts"] = time.Now().Format(val.(string))
	}
	for _, modreg := range Device.Modregs {
		switch modreg.FunctCode {
		case 1:
			results, err = r.RTUClient.ReadCoils(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)

		case 2:
			results, err = r.RTUClient.ReadDiscreteInputs(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		case 3:
			results, err = r.RTUClient.ReadHoldingRegisters(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		case 4:
			results, err = r.RTUClient.ReadInputRegisters(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		}
	}
	if val, ok := AddFeatures["ets"]; ok {
		tpay["ets"] = time.Now().Format(val.(string))
	}
	tpay["duid"] = Device.Device.DeviceId
	ModWg.Wait()
	if err == nil {

		Payload <- tpay
	}

	close(Payload)
	MqWg.Wait() // remove MqWait once archive to DB is done in MQTT side
}

func (t *TcpContainer) ModTrigger(Device models.DeviceWithRegs, ctx context.Context, TcpDevWg *sync.WaitGroup, DaqRate int) {

	for {
		select {
		case <-ctx.Done():
			TcpDevWg.Done()
			return
		case <-time.After(time.Duration(DaqRate) * time.Millisecond):
			go t.ModOperationAsync(Device, ctx, TcpDevWg)
		}
	}
}

func (t *TcpContainer) ModOperationAsync(Device models.DeviceWithRegs, ctx context.Context, TcpDevWg *sync.WaitGroup) {
	Payload := make(chan map[string]interface{}, len(Device.Modregs))
	tpay := make(map[string]interface{})
	var ModWg sync.WaitGroup
	var MqWg sync.WaitGroup
	MqWg.Add(1)                                       // remove MqWait once archive to DB is done in MQTT side
	go mq.MqttPublish2(Device.Device, Payload, &MqWg) // change to go mq.MqttPublish2(Device.Device, Payload)" once archive to DB is done in MQTT side
	t.TCPClientHandler.SlaveId = Device.Device.MbId

	if val, ok := AddFeatures["sts"]; ok {
		tpay["sts"] = time.Now().Format(val.(string))
	}

	for _, modreg := range Device.Modregs {
		switch modreg.FunctCode {
		case 1:
			results, err := t.TCPClient.ReadCoils(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		case 2:
			results, err := t.TCPClient.ReadDiscreteInputs(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		case 3:
			results, err := t.TCPClient.ReadHoldingRegisters(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		case 4:
			results, err := t.TCPClient.ReadInputRegisters(modreg.Register, modreg.Qty)
			ModWg.Add(1)
			go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
		}
	}
	tpay["duid"] = Device.Device.DeviceId

	ModWg.Wait()
	if val, ok := AddFeatures["ets"]; ok {
		tpay["ets"] = time.Now().Format(val.(string))
	}
	Payload <- tpay
	close(Payload)
	MqWg.Wait() //remove MqWait once archive to DB is done in MQTT side
}

func (t *TcpContainer) ModOperationSync(Device models.DeviceWithRegs, ctx context.Context, TcpDevWg *sync.WaitGroup, DaqRate int) {
	var ModWg sync.WaitGroup
	var MqWg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			TcpDevWg.Done()
			return
		case <-time.After(time.Duration(DaqRate) * time.Millisecond):
			Payload := make(chan map[string]interface{}, len(Device.Modregs))
			tpay := make(map[string]interface{})
			MqWg.Add(1)                                       // remove MqWait once archive to DB is done in MQTT side
			go mq.MqttPublish2(Device.Device, Payload, &MqWg) // change to go mq.MqttPublish2(Device.Device, Payload)" once archive to DB is done in MQTT side
			t.TCPClientHandler.SlaveId = Device.Device.MbId
			if val, ok := AddFeatures["sts"]; ok {
				tpay["sts"] = time.Now().Format(val.(string))
			}
			for _, modreg := range Device.Modregs {
				switch modreg.FunctCode {
				case 1:
					results, err := t.TCPClient.ReadCoils(modreg.Register, modreg.Qty)
					ModWg.Add(1)
					go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
				case 2:
					results, err := t.TCPClient.ReadDiscreteInputs(modreg.Register, modreg.Qty)
					ModWg.Add(1)
					go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
				case 3:
					results, err := t.TCPClient.ReadHoldingRegisters(modreg.Register, modreg.Qty)
					ModWg.Add(1)
					go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
				case 4:
					results, err := t.TCPClient.ReadInputRegisters(modreg.Register, modreg.Qty)
					ModWg.Add(1)
					go ModReadDataProcess(Device.Device, modreg, results, Payload, &ModWg, err)
				}
			}
			tpay["duid"] = Device.Device.DeviceId

			ModWg.Wait()
			if val, ok := AddFeatures["ets"]; ok {
				tpay["ets"] = time.Now().Format(val.(string))
			}
			Payload <- tpay
			close(Payload)
			MqWg.Wait()
		}
	}
}
