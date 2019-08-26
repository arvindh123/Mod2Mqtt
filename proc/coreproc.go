package proc

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/arvindh123/Mod2Mqtt/addparams"
	"github.com/arvindh123/Mod2Mqtt/controllers"
	"github.com/arvindh123/Mod2Mqtt/genrand"
	"github.com/arvindh123/Mod2Mqtt/mod"
	"github.com/arvindh123/Mod2Mqtt/mq"
	"github.com/arvindh123/Mod2Mqtt/ws"
)

var modStatus, mqStatus bool = false, false
var mqErr error
var processMu sync.Mutex
var ctx context.Context
var cancel context.CancelFunc
var AddFeatures map[string]interface{}
var autoStartMem bool = false
var autoStart bool = false
var auotStartTime uint64 = 1000

func ManInit() {
	AddFeatures = addparams.GetParams()
	if val, ok := AddFeatures["AutoStart"]; ok {
		autoStart = val.(bool)
	}
	if val, ok := AddFeatures["AutoStartTime"]; ok {
		auotStartTime = val.(uint64)
	}

	if autoStart {
		autoStartMem = true
	} else {
		autoStartMem = false
	}

}

func TryToStart(StartChan chan bool, StopChan chan bool) {
	for {

		select {
		case <-time.After(time.Duration(auotStartTime) * time.Millisecond):
			if !(modStatus && mqStatus) && autoStartMem {
				StartProc()
			}
		case cmd := <-StartChan:
			if cmd {
				ManInit()
				StartProc()
			}
		case cmd := <-StopChan:
			if cmd {
				StopProc()
			}
		}
	}
}

func StartProc() {
	processMu.Lock()
	defer processMu.Unlock()
	var wg sync.WaitGroup
	if modStatus == false {
		wg.Add(1)
		go StartMod(AddFeatures, &wg)
	} else {
		go ws.WsStatusPub(fmt.Sprintf("%v	Modbus Process is already started", time.Now()))
	}

	if mqStatus == false {
		wg.Add(1)
		go StartMq(AddFeatures, &wg)
	} else {
		go ws.WsStatusPub(fmt.Sprintf("%v	Mqtt Process is already started", time.Now()))
	}
	if autoStart {
		autoStartMem = true
	} else {
		autoStartMem = false
	}
	wg.Wait()
	return
}

func StopProc() {
	go ws.WsStatusPub(fmt.Sprintf("%v	No of Gorotuines before stop- %d", time.Now(), runtime.NumGoroutine()))
	var wg sync.WaitGroup
	wg.Add(1)
	go StopMod(cancel, &wg)
	wg.Add(1)
	go StopMq(&wg)
	wg.Wait()
	go ws.WsStatusPub(fmt.Sprintf("%v	No of Gorotuines After stop- %d", time.Now(), runtime.NumGoroutine()))
	autoStartMem = false
	return
}

func StartMq(AddFeatures map[string]interface{}, wg *sync.WaitGroup) {
	// fmt.Println("At Start mod mq status -", modStatus, mqStatus)
	if mqStatus == false {
		mqParams, err := controllers.GetMq()
		if err == nil {
			if Cid, ok := AddFeatures["EdgeUID"]; ok {
				mqStatus, mqErr = mq.MqConnect(Cid.(string), mqParams)
			} else {
				Cid := "mod2mqtt_" + genrand.RandStringRunes(5) + "_"
				mqStatus, mqErr = mq.MqConnect(Cid, mqParams)
				fmt.Println("mqStatus recevied - ", mqStatus)
			}
		} else {
			go ws.WsStatusPub(fmt.Sprintf("%v	Error Fetching Mqtt Details ", time.Now(), err.Error()))
		}
	} else {
		go ws.WsStatusPub(fmt.Sprintf("%v	MQTT Process is already running", time.Now()))
	}
	wg.Done()
	// fmt.Println("At end mod mq status -", modStatus, mqStatus)
}
func StartMod(AddFeatures map[string]interface{}, wg *sync.WaitGroup) {
	if modStatus == false {
		Inters := controllers.GetAllInters()
		addToTopic := ""
		if status, ok := AddFeatures["AddRootTopic"]; ok {
			if status.(bool) == true {
				if val, ok := AddFeatures["RootTopic"]; ok {
					addToTopic = addToTopic + val.(string) + "/"
				}
			}
		}
		if status, ok := AddFeatures["AddEdgeUidToTopic"]; ok {
			if status.(bool) == true {
				if val, ok := AddFeatures["EdgeUID"]; ok {
					addToTopic = addToTopic + val.(string) + "/"
				}
			}
		}

		for i := range Inters.AllStructParams {
			for j := range Inters.AllStructParams[i].Devices {
				if len(addToTopic) > 0 {
					Inters.AllStructParams[i].Devices[j].Device.PubTopic = addToTopic + Inters.AllStructParams[i].Devices[j].Device.DeviceId + "/" + Inters.AllStructParams[i].Devices[j].Device.PubTopic
				} else {
					Inters.AllStructParams[i].Devices[j].Device.PubTopic = Inters.AllStructParams[i].Devices[j].Device.DeviceId + "/" + Inters.AllStructParams[i].Devices[j].Device.PubTopic
				}
			}
		}
		ctx, cancel = context.WithCancel(context.Background())
		modStatus = mod.ModProcStart(Inters, ctx, AddFeatures)

	} else {
		go ws.WsStatusPub(fmt.Sprintf("%v	Modbus Process is already running", time.Now()))
	}
	wg.Done()
}
func StopMq(wg *sync.WaitGroup) {
	mqStatus = mq.MqttProcStop()
	wg.Done()
	if mqStatus {
		go ws.WsStatusPub(fmt.Sprintf("%v	MQTT Process is not able to stop", time.Now()))
		return
	}
	go ws.WsStatusPub(fmt.Sprintf("%v	MQTT Process stopped", time.Now()))
	return
}
func StopMod(cancel context.CancelFunc, wg *sync.WaitGroup) {
	modStatus = mod.ModProcStop(cancel)
	wg.Done()
	if modStatus {
		go ws.WsStatusPub(fmt.Sprintf("%v	Modbus Process is not able to stop", time.Now()))
		return
	}
	go ws.WsStatusPub(fmt.Sprintf("%v	Modbus Process stopped", time.Now()))
	return
}
