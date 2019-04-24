package mod

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"../models"
	"github.com/PaesslerAG/gval"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/modbus"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

type IdStruct struct {
	Id int
}

type PortStruct struct {
	Id      int
	ComPort string
}

type ModSpawner struct {
	ModRegs    []models.ModbusRegisters
	Handler    *modbus.RTUClientHandler
	ComChan    chan map[string]interface{}
	PortChance chan int
}

func (mod *ModSpawner) init() {
	mod.ComChan = make(chan map[string]interface{}, 1)
	mod.PortChance = make(chan int, 1)
}

type StatusLog struct {
	StatusLog string `json:"statuslog"`
}
type ModLastAquired struct {
	ModLastAquired string `json:"modlastaquired"`
}
type MqLastSent struct {
	MqLastSent string `json:"mqlastsent"`
}

type JSONString string

func (j JSONString) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

var WsClients = make(map[*websocket.Conn]bool)
var mu, processMu sync.Mutex

func MultiModMqProcessStart(db *gorm.DB, status bool, wsClientsChan chan map[*websocket.Conn]bool) (bool, []*modbus.RTUClientHandler, mqtt.Client, map[string](chan int), map[string](chan int), *sync.WaitGroup, []error, error) {
	processMu.Lock()
	defer processMu.Unlock()

	var handlers []*modbus.RTUClientHandler
	var retErrs []error
	var wg sync.WaitGroup
	var mqClient mqtt.Client
	var err, mqErr error
	var count int = 0
	var Topics []models.Topics
	_ = Topics
	// var status bool
	// status = false

	PortChance := make(map[string](chan int))

	SpanStopper := make(map[string](chan int))

	mqClient, mqErr = MqConnect(db)
	if mqErr != nil {
		// fmt.Println("Error in Mqtt Start ", mqErr)
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Error in Mqtt Start - %s", mqErr.Error(), time.Now()))
		return status, handlers, mqClient, SpanStopper, PortChance, &wg, nil, mqErr
	} else {
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Mqtt Connected", time.Now()))
		handlers, retErrs = MultiModStart(db)
		for i, err := range retErrs {
			if err != nil {
				// fmt.Println("Error in Modstart ", err)
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Error in Modstart - %s,  Serial Port - %s ", time.Now(), err.Error(), handlers[i].Address))
				if i < cap(handlers) {
					handlers[i].Close()
				}
			} else {
				count = count + 1
				PortChance[handlers[i].Address] = make(chan int, 1)
			}
		}
		if count > 0 {
			Topics, err = GetAllTopics(db)
			if err == nil {
				if len(Topics) > 0 {
					for _, topic := range Topics {
						tempTopicReqPorts := GetTopicReqPorts(db, topic.Id)
						if len(tempTopicReqPorts) > 0 {
							wg.Add(1)
						}

						SpanStopper[topic.Topic] = make(chan int, 1)
					}

					go MultiModReadWrite(db, mqClient, Topics, handlers, SpanStopper, PortChance, &wg)
					status = true
				} else {
					// fmt.Println("There is no topics, Please add topics")
					go WsStatusPub(WsClients, fmt.Sprintf("%v	There is no topics, Please add topics", time.Now()))
				}
			} else {
				// fmt.Println("Error in Fetching topics ", err)
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Error in Fetching topics ", time.Now(), err.Error()))
			}
		}
		return status, handlers, mqClient, SpanStopper, PortChance, &wg, retErrs, mqErr
	}

}

func MultiModMqProcessStop(handlers []*modbus.RTUClientHandler, mqClient mqtt.Client, SpanStopper map[string](chan int), PortChance map[string](chan int), wg *sync.WaitGroup) error {
	processMu.Lock()
	defer processMu.Unlock()

	go WsStatusPub(WsClients, fmt.Sprintf("%v	No of Gorotuines - %d", time.Now(), runtime.NumGoroutine()))
	go WsStatusPub(WsClients, fmt.Sprintf("%v	Topic Spanner and Modbus Ports closing", time.Now()))
	// fmt.Println("Topic Spanner and Modbus Ports closing")

	for k, v := range SpanStopper {
		v <- 100
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Stop Spanning topic - %s", time.Now(), k))
		// fmt.Println("Stop Spanning topic - %s", k)
	}

	// for k, v := range PortChance {
	// 	v <- 0
	// 	fmt.Println("Stoping the Port %s", k)

	// }

	// fmt.Println("Waiting for WG")
	go WsStatusPub(WsClients, fmt.Sprintf("%v	Waiting for WG", time.Now()))
	wg.Wait()
	// fmt.Println("Finished all WG")
	go WsStatusPub(WsClients, fmt.Sprintf("%v	Finished all WG", time.Now()))
	for _, v := range SpanStopper {
		if len(v) > 0 {
			<-v
		}
	}
	for _, v := range PortChance {
		if len(v) == 0 {
			v <- 0
		}
	}
	for _, v := range PortChance {
		if len(v) > 0 {
			<-v
		}
	}
	for _, handler := range handlers {
		handler.Close()
	}

	mqClient.Disconnect(250)

	go WsStatusPub(WsClients, fmt.Sprintf("%v	Topic Spanner and Modbus Ports Closed", time.Now()))
	// fmt.Println("Topic Spanner and Modbus Ports Closed")
	// go WsStatusPub(WsClients, fmt.Sprintf("%v	No of Gorotuines - %d", time.Now(), runtime.NumGoroutine()))
	return nil
}

func GetAllTopics(db *gorm.DB) ([]models.Topics, error) {
	var topics []models.Topics
	if err := db.Find(&topics).Error; err != nil {
		return topics, err
	} else {
		for i, topic := range topics {
			var modregs []*models.ModbusRegisters
			if err = db.Model(&topic).Association("ModRegs").Find(&modregs).Error; err != nil {
				// topics[i].ModRegs = modregs
			} else {
				// for i, modreg := range modregs {
				// 	var ports []*models.SerialDetails
				// 	if err = db.Model(&modreg).Association("SerialPorts").Find(&ports).Error; err != nil {
				// 	} else {
				// 		modregs[i].SerialPorts = ports
				// 	}
				// }
				topics[i].ModRegs = modregs
			}
		}
		return topics, nil
	}
}
func GetReqPorts(db *gorm.DB) []int {
	var temp []IdStruct
	var id []int
	if err := db.Raw(`select distinct SD.id FROM topics TT 
					   INNER join regs_topics RT ON TT.id = RT.topics_id 
					   INNER join modbus_registers MR ON MR.id = RT.modbus_registers_id 
					   INNER Join regs_ports RP ON RP.modbus_registers_id = MR.id 
					   inner join serial_details SD ON SD.id = RP.serial_details_id`).Scan(&temp).Error; err != nil {
	}
	for _, st := range temp {
		id = append(id, st.Id)
	}
	return id
}

func getHandlerFromHandlers(handlers []*modbus.RTUClientHandler, PortName string) *modbus.RTUClientHandler {
	var Handler *modbus.RTUClientHandler
	for _, handler := range handlers {
		if PortName == handler.Address {
			Handler = handler
			break
		}
	}
	return Handler
}
func GetTopicPortReqModRegs(db *gorm.DB, TopId int, Port string) []models.ModbusRegisters {
	var temp []models.ModbusRegisters
	if err := db.Raw(`	select distinct MR.id, MR.name, MR.unit, MR.funct_code, MR.register, MR.qty, MR.data_type, 
						MR.byte_order, MR.post_process, MR.tags FROM topics TT 
						INNER join regs_topics RT ON TT.id = RT.topics_id 
						INNER join modbus_registers MR ON MR.id = RT.modbus_registers_id 
						INNER Join regs_ports  RP ON RP.modbus_registers_id = MR.id
						inner join serial_details SD ON SD.id = RP.serial_details_id
						where TT.id = ? AND SD.com_port = ?`, TopId, Port).Scan(&temp).Error; err != nil {
	}

	return temp

}

func GetTopicReqPorts(db *gorm.DB, ReqId int) []PortStruct {
	var temp []PortStruct
	if err := db.Raw(`select distinct SD.id, SD.com_port FROM topics TT 
					   INNER join regs_topics RT ON TT.id = RT.topics_id 
					   INNER join modbus_registers MR ON MR.id = RT.modbus_registers_id 
					   INNER Join regs_ports RP ON RP.modbus_registers_id = MR.id 
					   inner join serial_details SD ON SD.id = RP.serial_details_id 
					   where TT.id = ?`, ReqId).Scan(&temp).Error; err != nil {
	}

	return temp

}

func GetReqPortsOld(db *gorm.DB) []int {
	topics, _ := GetAllTopics(db)
	// var topics1 models.Topics
	var id []int

	for _, topic := range topics {
		for _, modreg := range topic.ModRegs {
			for _, port := range modreg.SerialPorts {
				id = append(id, port.Id)
			}
		}
	}
	return id
}

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

func WsStatusPub(clients map[*websocket.Conn]bool, msg string) {
	mu.Lock()
	defer mu.Unlock()
	wsJson, err := json.Marshal(StatusLog{msg})
	if err == nil {
		for client := range WsClients {
			err := client.WriteMessage(1, wsJson)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func MultiModStart(db *gorm.DB) ([]*modbus.RTUClientHandler, []error) {
	var serialParams []models.SerialDetails
	var handlers []*modbus.RTUClientHandler
	var retErrs []error
	if err := db.Where(GetReqPorts(db)).Find(&serialParams).Error; err != nil {
		retErrs = append(retErrs, err)
		return handlers, retErrs
	} else {
		for _, serialParam := range serialParams {
			handler := modbus.NewRTUClientHandler(serialParam.ComPort)
			handler.BaudRate = serialParam.BaudRate
			handler.DataBits = serialParam.DataBits
			handler.Parity = serialParam.Parity
			handler.StopBits = serialParam.StopBits
			handler.Timeout = time.Duration(serialParam.Timeout) * time.Second
			// handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
			err := handler.Connect()
			handlers = append(handlers, handler)
			retErrs = append(retErrs, err)
		}
		return handlers, retErrs
	}
}

func ModReconn(handler *modbus.RTUClientHandler) (*modbus.RTUClientHandler, error) {
	for {
		err := handler.Connect()
		if err == nil {
			return handler, err
		} else if err != nil {
			return handler, err
		}
	}
}

func MultiModReadWrite(db *gorm.DB, mqClient mqtt.Client, Topics []models.Topics, handlers []*modbus.RTUClientHandler, SpanStopper map[string](chan int), PortChance map[string](chan int), wg *sync.WaitGroup) {
	go WsStatusPub(WsClients, fmt.Sprintf("%v	No of Gorotuines - %d", time.Now(), runtime.NumGoroutine()))
	// fmt.Println("Modbus RW Started")
	go WsStatusPub(WsClients, fmt.Sprintf("%v	Modbus RW Started", time.Now()))

	for _, topic := range Topics {
		// go MqttOperation(mqClient, topic, payload, mqShutDown, wg2)
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Goroutine STARTED for Publishing Mqtt Topic - %s", time.Now(), topic.Topic))
		go TopicModPort(db, mqClient, topic, handlers, SpanStopper[topic.Topic], PortChance, wg)
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Goroutine STARTED for Reading Modbus Registers - %v which are related to Topic - %s", time.Now(), topic.ModRegs, topic.Topic))
	}

	for k, v := range PortChance {
		v <- 100
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Send start cmd 100 to Port - %s", time.Now(), k))
	}
	return
}

func TopicModPort(db *gorm.DB, mqClient mqtt.Client, topic models.Topics, handlers []*modbus.RTUClientHandler, SpanStopperTopic chan int, PortChance map[string](chan int), wg *sync.WaitGroup) {
	TopicReqPorts := GetTopicReqPorts(db, topic.Id)
	var RegsPorts []ModSpawner
	if len(TopicReqPorts) > 0 {

		for _, Port := range TopicReqPorts {
			h := getHandlerFromHandlers(handlers, Port.ComPort)
			RegsPort := ModSpawner{ModRegs: GetTopicPortReqModRegs(db, topic.Id, h.Address), Handler: h}
			RegsPort.init()
			RegsPort.PortChance = PortChance[h.Address]
			RegsPorts = append(RegsPorts, RegsPort)
		}

		go MultiModSpwaner2(mqClient, topic, SpanStopperTopic, RegsPorts, wg)
	}

	// return
}

// func MultiModSpwaner(mqClient mqtt.Client, topic models.Topics, SpanStopperTopic chan int, RegsPorts []ModSpawner) {
// 	var Pays []chan map[string]interface{}
// 	for _, RegsPort := range RegsPorts {
// 		Pays = append(Pays, RegsPort.ComChan)
// 	}
// 	for {
// 		Payload := make(map[string]interface{})
// 		_ = Payload
// 		var wg sync.WaitGroup
// 		wg.Add(len(RegsPorts))
// 		for _, RegsPort := range RegsPorts {
// 			go ModOperation(topic, RegsPort, &wg)
// 		}
// 		go MqttPublish(mqClient, topic, RegsPorts)
// 		wg.Wait()
// 		time.Sleep(time.Duration(topic.Delay) * time.Second)
// 	}
// }

// func ModOperation(topic models.Topics, RegsPort ModSpawner, wg *sync.WaitGroup) {
// 	client := modbus.NewClient(RegsPort.Handler)
// 	tpay := make(map[string]interface{})
// 	_ = tpay
// 	lenOfModregs := cap(RegsPort.ModRegs)
// 	_ = lenOfModregs
// 	for _, modreg := range RegsPort.ModRegs {
// 		time.Sleep(time.Duration(topic.Delay) * time.Second)
// 		RegsPort.Handler.SlaveId = modreg.Unit
// 		results, err := client.ReadHoldingRegisters(modreg.Register, modreg.Qty)
// 		Processed := ModReadDataProcess(topic, &modreg, results, err)

// 		if len(Processed) > 0 {

// 			for k, v := range Processed {
// 				tpay[k] = v
// 			}

// 		}
// 	}
// 	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
// 	RegsPort.ComChan <- tpay
// 	wg.Done()
// }

// func MqttPublish(mqClient mqtt.Client, topic models.Topics, RegsPorts []ModSpawner) {
// 	finalPay := make(map[string]interface{})
// 	for _, RegsPort := range RegsPorts {
// 		for k, v := range <-RegsPort.ComChan {
// 			finalPay[k] = v
// 		}
// 	}
// 	// fmt.Println("Pusblished data %v", finalPay)
// 	// mqClient.Publish(topic.Topic, 0, false, payload)
// 	// fmt.Sprintf("Topic - %s, Payload- %s", topic.Topic, payload)
// 	// go WsClientPub(MqLastSent{fmt.Sprintf("Topic - %s, Payload- %s", topic.Topic, payload)})
// }

func MultiModSpwaner2(mqClient mqtt.Client, topic models.Topics, SpanStopperTopic chan int, RegsPorts []ModSpawner, wgMain *sync.WaitGroup) {
	for {
		select {
		case msg := <-SpanStopperTopic:
			if msg == 100 {
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Stopping the Topic Spanner - %s", time.Now(), topic.Topic))
				// close(Payload)
				wgMain.Done()
				return
			}
		default:
			Payload := make(chan map[string]interface{}, len(RegsPorts))
			var ModWg sync.WaitGroup
			// var MqWg sync.WaitGroup
			ModWg.Add(len(RegsPorts))
			// MqWg.Add(1)

			for _, RegsPort := range RegsPorts {
				RegsPort.ComChan = make(chan map[string]interface{})
				go ModOperation2(topic, RegsPort, Payload, &ModWg)
				// fmt.Println("Spanned")
			}
			// go MqttPublish2(mqClient, topic, Payload, &MqWg)
			go MqttPublish2(mqClient, topic, Payload)
			// fmt.Println("Waiting for wait group modbus reder in Spwaner")
			ModWg.Wait()
			close(Payload)
			// MqWg.Wait()
		}

		if topic.Delay > 0 {
			time.Sleep(time.Duration(topic.Delay) * time.Second)

		}

	}
	// return

}

func ModOperation2(topic models.Topics, RegsPort ModSpawner, Payload chan map[string]interface{}, wg *sync.WaitGroup) {
	select {
	case msg := <-RegsPort.PortChance:
		// fmt.Println("Got Msg - %d in port - %s of topic - %s", msg, RegsPort.Handler.Address, topic.Topic)
		if msg == 100 {
			client := modbus.NewClient(RegsPort.Handler)
			tpay := make(map[string]interface{})
			_ = tpay
			lenOfModregs := cap(RegsPort.ModRegs)
			_ = lenOfModregs
			for _, modreg := range RegsPort.ModRegs {
				time.Sleep(time.Duration(topic.Delay) * time.Second)
				RegsPort.Handler.SlaveId = modreg.Unit
				results := []byte{}
				err := errors.New("")
				switch modreg.FunctCode {
				case 1:
					results, err = client.ReadCoils(modreg.Register, modreg.Qty)
				case 2:
					results, err = client.ReadDiscreteInputs(modreg.Register, modreg.Qty)
				case 3:
					results, err = client.ReadHoldingRegisters(modreg.Register, modreg.Qty)
				case 4:
					results, err = client.ReadInputRegisters(modreg.Register, modreg.Qty)
				case 90:
					Idd, err := strconv.Atoi(modreg.PostProcess)
					if err == nil {
						tpay["suid"] = Idd

					} else {
						tpay["suid"] = 99999
					}
				case 99:
					tpay["ts"] = int32(time.Now().Unix())
				case 98:
					tpay["ts"] = time.Now().Round(0)
				}
				if modreg.FunctCode < 90 {
					Processed := ModReadDataProcess(topic, &modreg, results, err)
					if len(Processed) > 0 {
						for k, v := range Processed {
							tpay[k] = v
						}
					}
				}

			}

			Payload <- tpay

			if len(RegsPort.PortChance) == 0 {
				RegsPort.PortChance <- 100
			}
			wg.Done()

			return
		} else {
			if len(RegsPort.PortChance) == 0 {
				RegsPort.PortChance <- 0
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Stopping Mod Operation of the Topic Spanner - %s, Port - %s", time.Now(), topic.Topic, RegsPort.Handler.Address))
			}
			wg.Done()
			return

		}
	}

}

func MqttPublish2(mqClient mqtt.Client, topic models.Topics, Payload chan map[string]interface{}) { //, wg *sync.WaitGroup) {
	finalPay := make(map[string]interface{})
	// fmt.Println("Starting <- Pays")
	for vs := range Payload {
		// fmt.Println("for each payload - ", vs)
		for k, v := range vs {
			finalPay[k] = v
		}
	}
	final, err := json.Marshal(finalPay)
	// fmt.Println("Pusblished data %v", finalPay)
	if err == nil {
		mqClient.Publish(topic.Topic, 0, false, final)
		go WsClientPub(MqLastSent{fmt.Sprintf("Topic - %s, Payload- %v", topic.Topic, finalPay)})
	} else {
		go WsClientPub(MqLastSent{fmt.Sprintf("Topic - %s, Before Marshal Payload- %v, Error in Marshal -> Err- %v", topic.Topic, finalPay, err)})
	}

	// wg.Done()

	// return
}

func ModReadDataProcess(topic models.Topics, modreg *models.ModbusRegisters, results []byte, err error) map[string]interface{} {
	retString := make(map[string]interface{})
	if err == nil {
		if modreg.DataType == 1 {
			// 1 - uint 8 // 3 - int 8
			val := results
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}
		} else if modreg.DataType == 2 {
			// []uint 8
			val := results
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				if strings.Contains(modreg.Tags, ",") {
					AllTags := strings.Split(modreg.Tags, ",")
					for i, Tag := range AllTags {
						if i == (len(AllTags)-1) && len(AllTags) < len(val) {
							retString[Tag] = val[i:]
						} else {
							retString[Tag] = val[i]
						}

					}
				} else {
					retString[modreg.Tags] = val
				}
			}
		} else if modreg.DataType == 3 {
			// int 8
		} else if modreg.DataType == 4 {
			// []int 8
		} else if modreg.DataType == 5 {
			// uint16
			if strings.Contains(modreg.PostProcess, "value") {
				val, err := SingleUint16FromBytesPP(results, modreg.ByteOrder, modreg.PostProcess)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			} else {
				val, err := SingleUint16FromBytes(results, modreg.ByteOrder)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			}

		} else if modreg.DataType == 6 {
			// []uint16
			if strings.Contains(modreg.PostProcess, "value") {
				val, err := arrUint16frombytesPP(results, modreg.ByteOrder, modreg.PostProcess)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					if strings.Contains(modreg.Tags, ",") {
						AllTags := strings.Split(modreg.Tags, ",")
						for i, Tag := range AllTags {
							if i == (len(AllTags)-1) && len(AllTags) < len(val) {
								retString[Tag] = val[i:]
							} else {
								retString[Tag] = val[i]
							}

						}
					} else {
						retString[modreg.Tags] = val
					}
				}
			} else {
				val, err := arrUint16frombytes(results, modreg.ByteOrder)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					if strings.Contains(modreg.Tags, ",") {
						AllTags := strings.Split(modreg.Tags, ",")
						for i, Tag := range AllTags {
							if i == (len(AllTags)-1) && len(AllTags) < len(val) {
								retString[Tag] = val[i:]
							} else {
								retString[Tag] = val[i]
							}

						}
					} else {
						retString[modreg.Tags] = val
					}
				}
			}

		} else if modreg.DataType == 7 {
			// int16
			if strings.Contains(modreg.PostProcess, "value") {
				val, err := SingleUint16FromBytesPP(results, modreg.ByteOrder, modreg.PostProcess)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			} else {
				val, err := SingleUint16FromBytes(results, modreg.ByteOrder)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			}

		} else if modreg.DataType == 8 {
			// []int16
			if strings.Contains(modreg.PostProcess, "value") {
				val, err := arrUint16frombytesPP(results, modreg.ByteOrder, modreg.PostProcess)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					if strings.Contains(modreg.Tags, ",") {
						AllTags := strings.Split(modreg.Tags, ",")
						for i, Tag := range AllTags {
							if i == (len(AllTags)-1) && len(AllTags) < len(val) {
								retString[Tag] = val[i:]
							} else {
								retString[Tag] = val[i]
							}

						}
					} else {
						retString[modreg.Tags] = val
					}
				}
			} else {
				val, err := arrUint16frombytes(results, modreg.ByteOrder)
				go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					if strings.Contains(modreg.Tags, ",") {
						AllTags := strings.Split(modreg.Tags, ",")
						for i, Tag := range AllTags {
							if i == (len(AllTags)-1) && len(AllTags) < len(val) {
								retString[Tag] = val[i:]
							} else {
								retString[Tag] = val[i]
							}

						}
					} else {
						retString[modreg.Tags] = val
					}
				}
			}
		} else if modreg.DataType == 9 {
			// uint32
		} else if modreg.DataType == 10 {
			// []uint32
		} else if modreg.DataType == 11 {
			// int32
			val, err := SingleUint32FromBytes(results, modreg.ByteOrder)
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}
		} else if modreg.DataType == 12 {
			// []int32
			val, err := arrUint32frombytes(results, modreg.ByteOrder)
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				if strings.Contains(modreg.Tags, ",") {
					AllTags := strings.Split(modreg.Tags, ",")
					for i, Tag := range AllTags {
						if i == (len(AllTags)-1) && len(AllTags) < len(val) {
							retString[Tag] = val[i:]
						} else {
							retString[Tag] = val[i]
						}

					}
				} else {
					retString[modreg.Tags] = val
				}
			}
		} else if modreg.DataType == 13 {
			// uint64
		} else if modreg.DataType == 14 {
			// []uint64
		} else if modreg.DataType == 15 {
			// int64
			val, err := SingleUint64FromBytes(results, modreg.ByteOrder)
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}
		} else if modreg.DataType == 16 {
			// []int64
			val, err := arrUint64frombytes(results, modreg.ByteOrder)
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				if strings.Contains(modreg.Tags, ",") {
					AllTags := strings.Split(modreg.Tags, ",")
					for i, Tag := range AllTags {
						if i == (len(AllTags)-1) && len(AllTags) < len(val) {
							retString[Tag] = val[i:]
						} else {
							retString[Tag] = val[i]
						}

					}
				} else {
					retString[modreg.Tags] = val
				}
			}
		} else if modreg.DataType == 17 {
			// float32
			val, err := SingleFloat32FromBytes(results, modreg.ByteOrder)
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %f, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}

		} else if modreg.DataType == 18 {
			// []float32
			val, err := arrFloat32frombytes(results, modreg.ByteOrder)
			go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %f, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				if strings.Contains(modreg.Tags, ",") {
					AllTags := strings.Split(modreg.Tags, ",")
					for i, Tag := range AllTags {
						if i == (len(AllTags)-1) && len(AllTags) < len(val) {
							retString[Tag] = val[i:]
						} else {
							retString[Tag] = val[i]
						}

					}
				} else {
					retString[modreg.Tags] = val
				}
			}

		} else if modreg.DataType == 19 {
			// float64

		} else if modreg.DataType == 20 {
			// []float64

		} else if modreg.DataType == 21 {

		} else if modreg.DataType == 22 {

		} else if modreg.DataType == 23 {

		} else if modreg.DataType == 24 {

		} else {

		}
	} else {
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Topic- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", time.Now(), topic.Topic, modreg.Name, modreg.Register, results, err))
		// fmt.Println(topic.Topic, modreg.Name, modreg.Register, results, err)

	}
	return retString
}

func MqConnect(db *gorm.DB) (mqtt.Client, error) {
	var mqParams models.Mqtt
	if err := db.Where("id = ?", 1).First(&mqParams).Error; err != nil {
		return nil, err
	} else {
		clientid := fmt.Sprintf("TestClient_%d", rand.Intn(1000))
		server := fmt.Sprintf("tcp://"+mqParams.Ip+":%d", mqParams.Port)
		username := mqParams.Username
		password := mqParams.Password
		keepalive := 60
		opts := CreateMqClientOptions(clientid, server, username, password, keepalive)
		client := mqtt.NewClient(opts)
		token := client.Connect()
		for !token.WaitTimeout(3 * time.Second) {
		}
		if err := token.Error(); err != nil {
			log.Fatal(err)
		}
		return client, err
	}
}

func CreateMqClientOptions(clientid string, server string, username string, password string, keepalive int) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientid)
	// opts.SetConnectTimeout(time.Duration(5))
	// opts.SetKeepAlive(time.Duration(60))
	return opts
}

// func FloatValueToMap(a string , val []float32 ) map[string]float32 {

// 		arr := []int{1,2}

// 		s := strings.Split(a, ",")
// 		//fmt.Println(s)
// 		m := make(map[string]int)
// 		for idx, val := range s {
// 			if idx < len(arr) {
// 			m[val] = arr[idx]
// 			}else{

// 			}
// 		}
// 		fmt.Println(m)
// 	}

// func ArrayToString(A interface{}, delim string) string {
// 	typeoff := reflect.TypeOf(A)
// 	var buffer bytes.Buffer
// 	for i := 0; i < len(A); i++ {
// 		switch typeoff{
// 		case []int :
// 			buffer.WriteString(strconv.Itoa(A[i]))
// 		case int :
// 			buffer.WriteString(strconv.Itoa(A[i]))
// 		}

// 		if i != len(A)-1 {
// 			buffer.WriteString(delim)
// 		}
// 	}

// 	return buffer.String()
// }

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func SinglePostProcess(ppstring string, value interface{}) (interface{}, error) {
	val, err := gval.Evaluate(ppstring, map[string]interface{}{"value": value})
	if err != nil {
		go WsStatusPub(WsClients, fmt.Sprintf("Error in Post Process Ppstring -%s Value -%v", ppstring, value))
	}
	return val, err
}

func ArrPostProcess(ppstring string, values ...interface{}) {
	for _, value := range values {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			_ = s
		case reflect.Array:
			s := reflect.ValueOf(value)
			_ = s
		default:
			s := reflect.ValueOf(value)
			_ = s
		}
	}
}

func SingleFloat32FromBytes(bytes []byte, byteorder uint8) (float32, error) {
	bytes_len := len(bytes)
	var float float32
	var bits uint32
	if bytes_len == 4 {
		if byteorder == 1 { // comparison  1 = Big Endian
			bits = binary.BigEndian.Uint32(bytes)
			float = math.Float32frombits(bits)
			return float, nil
		} else if byteorder == 2 {
			bits = binary.LittleEndian.Uint32(bytes)
			float = math.Float32frombits(bits)
			return float, nil
		} else if byteorder == 3 {
			reordered := []byte{bytes[1], bytes[0], bytes[3], bytes[2]}
			bits = binary.BigEndian.Uint32(reordered)
			float = math.Float32frombits(bits)
			return float, nil
		} else if byteorder == 4 {
			reordered := []byte{bytes[1], bytes[0], bytes[3], bytes[2]}
			bits = binary.LittleEndian.Uint32(reordered)
			float = math.Float32frombits(bits)
			return float, nil
		} else {
			return 0, errors.New("Byte Order not specified")
		}
	} else {
		return 0, errors.New("Array length is not equal to 4")
	}
}

func arrFloat32frombytes(bytes []byte, byteorder uint8) ([]float32, error) {
	var splitnos, split_lim, bytes_len int
	splitnos = 4 //change this to 8 if the float64
	bytes_len = len(bytes)
	split_lim = bytes_len / splitnos
	var ret_float []float32
	var bits uint32
	var float float32

	if bytes_len >= 4 {
		if byteorder == 1 { // comparison  1 = Big Endian
			for i := 0; i < split_lim; i++ {
				bits = binary.BigEndian.Uint32(bytes[(i * splitnos):((i + 1) * splitnos)])
				float = math.Float32frombits(bits)

				ret_float = append(ret_float, float)
			}
		} else if byteorder == 2 { // comparison  2 = Little Endian
			for i := 0; i < split_lim; i++ {
				bits = binary.LittleEndian.Uint32(bytes[(i * splitnos):((i + 1) * splitnos)])
				float = math.Float32frombits(bits)
				ret_float = append(ret_float, float)
			}

		} else if byteorder == 3 { // comparison  3 = Mid-Big Endian
			for i := 0; i < split_lim; i++ {
				tempBytes := bytes[(i * splitnos):((i + 1) * splitnos)]
				reordered := []byte{tempBytes[1], tempBytes[0], tempBytes[3], tempBytes[2]}
				bits = binary.BigEndian.Uint32(reordered)
				float = math.Float32frombits(bits)
				ret_float = append(ret_float, float)
			}

		} else if byteorder == 4 { // comparison  4 = Mid-Littleittle Endian
			for i := 0; i < split_lim; i++ {
				tempBytes := bytes[(i * splitnos):((i + 1) * splitnos)]
				reordered := []byte{tempBytes[1], tempBytes[0], tempBytes[3], tempBytes[2]}
				bits = binary.LittleEndian.Uint32(reordered)
				float = math.Float32frombits(bits)
				ret_float = append(ret_float, float)
			}
		} else {
			ret_float = append(ret_float, 0)
			return ret_float, errors.New("Byte Order not specified")
		}
		return ret_float, nil
	} else {
		ret_float = append(ret_float, 0)
		return ret_float, errors.New("Array length is not greater than 4")
	}
}

func SingleUint32FromBytes(bytes []byte, byteorder uint8) (uint32, error) {
	bytes_len := len(bytes)
	var val uint32
	if bytes_len == 4 {
		if byteorder == 1 { // comparison  1 = Big Endian
			val = binary.BigEndian.Uint32(bytes)
			return val, nil
		} else if byteorder == 2 {
			val = binary.LittleEndian.Uint32(bytes)
			return val, nil
		} else if byteorder == 3 {
			reordered := []byte{bytes[1], bytes[0], bytes[3], bytes[2]}
			val = binary.BigEndian.Uint32(reordered)
			return val, nil
		} else if byteorder == 4 {
			reordered := []byte{bytes[1], bytes[0], bytes[3], bytes[2]}
			val = binary.LittleEndian.Uint32(reordered)
			return val, nil
		} else {
			return 0, errors.New("Byte Order not specified")
		}
	} else {
		return 0, errors.New("Array length is not equal to 4")
	}
}

func arrUint32frombytes(bytes []byte, byteorder uint8) ([]uint32, error) {
	var splitnos, split_lim, bytes_len int
	splitnos = 4 //change this to 8 if the float64
	bytes_len = len(bytes)
	split_lim = bytes_len / splitnos
	var ret_val []uint32
	var val uint32

	if bytes_len >= 4 {
		if byteorder == 1 { // comparison  1 = Big Endian
			for i := 0; i < split_lim; i++ {
				val = binary.BigEndian.Uint32(bytes[(i * splitnos):((i + 1) * splitnos)])
				ret_val = append(ret_val, val)
			}
		} else if byteorder == 2 { // comparison  2 = Little Endian
			for i := 0; i < split_lim; i++ {
				val = binary.LittleEndian.Uint32(bytes[(i * splitnos):((i + 1) * splitnos)])
				ret_val = append(ret_val, val)
			}

		} else if byteorder == 3 { // comparison  3 = Mid-Big Endian
			for i := 0; i < split_lim; i++ {
				tempBytes := bytes[(i * splitnos):((i + 1) * splitnos)]
				reordered := []byte{tempBytes[1], tempBytes[0], tempBytes[3], tempBytes[2]}
				val = binary.BigEndian.Uint32(reordered)
				ret_val = append(ret_val, val)
			}

		} else if byteorder == 4 { // comparison  4 = Mid-Littleittle Endian
			for i := 0; i < split_lim; i++ {
				tempBytes := bytes[(i * splitnos):((i + 1) * splitnos)]
				reordered := []byte{tempBytes[1], tempBytes[0], tempBytes[3], tempBytes[2]}
				val = binary.LittleEndian.Uint32(reordered)
				ret_val = append(ret_val, val)
			}
		} else {
			ret_val = append(ret_val, 0)
			return ret_val, errors.New("Byte Order not specified")
		}
		return ret_val, nil
	} else {
		ret_val = append(ret_val, 0)
		return ret_val, errors.New("Array length is not greater than 4")
	}
}

func SingleUint64FromBytes(bytes []byte, byteorder uint8) (uint64, error) {
	bytes_len := len(bytes)
	var val uint64
	if bytes_len == 8 {
		if byteorder == 1 { // comparison  1 = Big Endian
			val = binary.BigEndian.Uint64(bytes)
			return val, nil
		} else if byteorder == 2 {
			val = binary.LittleEndian.Uint64(bytes)
			return val, nil
		} else if byteorder == 3 {
			reordered := []byte{bytes[1], bytes[0], bytes[3], bytes[2], bytes[5], bytes[4], bytes[7], bytes[6]}
			val = binary.BigEndian.Uint64(reordered)
			return val, nil
		} else if byteorder == 4 {
			reordered := []byte{bytes[1], bytes[0], bytes[3], bytes[2], bytes[5], bytes[4], bytes[7], bytes[6]}
			val = binary.LittleEndian.Uint64(reordered)
			return val, nil
		} else {
			return 0, errors.New("Byte Order not specified")
		}
	} else {
		return 0, errors.New("Array length is not equal to 8")
	}
}

func arrUint64frombytes(bytes []byte, byteorder uint8) ([]uint64, error) {
	var splitnos, split_lim, bytes_len int
	splitnos = 8 //change this to 8 if the float64
	bytes_len = len(bytes)
	split_lim = bytes_len / splitnos
	var ret_val []uint64
	var val uint64

	if bytes_len >= 8 {
		if byteorder == 1 { // comparison  1 = Big Endian
			for i := 0; i < split_lim; i++ {
				val = binary.BigEndian.Uint64(bytes[(i * splitnos):((i + 1) * splitnos)])
				ret_val = append(ret_val, val)
			}
		} else if byteorder == 2 { // comparison  2 = Little Endian
			for i := 0; i < split_lim; i++ {
				val = binary.LittleEndian.Uint64(bytes[(i * splitnos):((i + 1) * splitnos)])
				ret_val = append(ret_val, val)
			}

		} else if byteorder == 3 { // comparison  3 = Mid-Big Endian
			for i := 0; i < split_lim; i++ {
				tempBytes := bytes[(i * splitnos):((i + 1) * splitnos)]
				reordered := []byte{tempBytes[1], tempBytes[0], tempBytes[3], tempBytes[2], tempBytes[5], tempBytes[4], tempBytes[7], tempBytes[6]}
				val = binary.BigEndian.Uint64(reordered)
				ret_val = append(ret_val, val)
			}

		} else if byteorder == 4 { // comparison  4 = Mid-Littleittle Endian
			for i := 0; i < split_lim; i++ {
				tempBytes := bytes[(i * splitnos):((i + 1) * splitnos)]
				reordered := []byte{tempBytes[1], tempBytes[0], tempBytes[3], tempBytes[2], tempBytes[5], tempBytes[4], tempBytes[7], tempBytes[6]}
				val = binary.LittleEndian.Uint64(reordered)
				ret_val = append(ret_val, val)
			}
		} else {
			ret_val = append(ret_val, 0)
			return ret_val, errors.New("Byte Order not specified")
		}
		return ret_val, nil
	} else {
		ret_val = append(ret_val, 0)
		return ret_val, errors.New("Array length is not greater than 8")
	}
}

func SingleUint16FromBytesPP(bytes []byte, byteorder uint8, ppstring string) (interface{}, error) {
	ppstrings := strings.Split(ppstring, "::")

	bytes_len := len(bytes)

	if bytes_len == 2 {
		if byteorder == 1 { // comparison  1 = Big Endian
			val, err := SinglePostProcess(ppstrings[0], binary.BigEndian.Uint16(bytes))
			return val, err
		} else if byteorder == 2 {
			val, err := SinglePostProcess(ppstrings[0], binary.LittleEndian.Uint16(bytes))
			return val, err
		} else {
			return 0, errors.New("Byte Order not specified")
		}
	} else {
		return 0, errors.New("Array length is not equal to 2")
	}
}

func arrUint16frombytesPP(bytes []byte, byteorder uint8, ppstring string) ([]interface{}, error) {
	ppstrings := strings.Split(ppstring, "::")
	capPpstring := cap(ppstrings)
	var splitnos, split_lim, bytes_len int
	splitnos = 2 //change this to 8 if the float64
	bytes_len = len(bytes)
	split_lim = bytes_len / splitnos
	var ret_val []interface{}
	var val interface{}
	var err error
	_ = err
	if bytes_len >= 4 {
		if byteorder == 1 { // comparison  1 = Big Endian
			for i := 0; i < split_lim; i++ {
				if i < capPpstring {
					val, err = SinglePostProcess(ppstrings[i], binary.BigEndian.Uint16(bytes[(i*splitnos):((i+1)*splitnos)]))
				} else {
					val, err = SinglePostProcess(ppstrings[0], binary.BigEndian.Uint16(bytes[(i*splitnos):((i+1)*splitnos)]))
				}
				ret_val = append(ret_val, val)
			}
		} else if byteorder == 2 { // comparison  2 = Little Endian
			for i := 0; i < split_lim; i++ {
				if i < capPpstring {
					val, err = SinglePostProcess(ppstrings[i], binary.LittleEndian.Uint16(bytes[(i*splitnos):((i+1)*splitnos)]))
				} else {
					val, err = SinglePostProcess(ppstrings[0], binary.LittleEndian.Uint16(bytes[(i*splitnos):((i+1)*splitnos)]))
				}
				ret_val = append(ret_val, val)
			}
		} else {
			ret_val = append(ret_val, 0)
			return ret_val, errors.New("Byte Order not specified")
		}
		return ret_val, nil
	} else {
		ret_val = append(ret_val, 0)
		return ret_val, errors.New("Array length is not greater than 2")
	}
}

func SingleUint16FromBytes(bytes []byte, byteorder uint8) (uint16, error) {

	bytes_len := len(bytes)
	var val uint16
	if bytes_len == 2 {
		if byteorder == 1 { // comparison  1 = Big Endian
			val = binary.BigEndian.Uint16(bytes)
			return val, nil
		} else if byteorder == 2 {
			val = binary.LittleEndian.Uint16(bytes)
			return val, nil
		} else {
			return 0, errors.New("Byte Order not specified")
		}
	} else {
		return 0, errors.New("Array length is not equal to 2")
	}
}

func arrUint16frombytes(bytes []byte, byteorder uint8) ([]uint16, error) {
	var splitnos, split_lim, bytes_len int
	splitnos = 2 //change this to 8 if the float64
	bytes_len = len(bytes)
	split_lim = bytes_len / splitnos
	var ret_val []uint16
	var val uint16

	if bytes_len >= 4 {
		if byteorder == 1 { // comparison  1 = Big Endian
			for i := 0; i < split_lim; i++ {
				val = binary.BigEndian.Uint16(bytes[(i * splitnos):((i + 1) * splitnos)])
				ret_val = append(ret_val, val)
			}
		} else if byteorder == 2 { // comparison  2 = Little Endian
			for i := 0; i < split_lim; i++ {
				val = binary.LittleEndian.Uint16(bytes[(i * splitnos):((i + 1) * splitnos)])
				ret_val = append(ret_val, val)
			}
		} else {
			ret_val = append(ret_val, 0)
			return ret_val, errors.New("Byte Order not specified")
		}
		return ret_val, nil
	} else {
		ret_val = append(ret_val, 0)
		return ret_val, errors.New("Array length is not greater than 2")
	}
}
