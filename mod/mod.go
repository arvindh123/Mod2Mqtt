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
	"sync"
	"time"

	"../models"
	"github.com/PaesslerAG/gval"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/goburrow/modbus"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
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

type JSONString string

func (j JSONString) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

var WsClients = make(map[*websocket.Conn]bool)
var mu, processMu sync.Mutex

func MultiModMqProcessStart(db *gorm.DB, wsClientsChan chan map[*websocket.Conn]bool) ([]*modbus.RTUClientHandler, []modbus.Client, mqtt.Client, *sync.WaitGroup, chan byte, error, []error) {
	// processMu.Lock()
	// defer processMu.Unlock()
	// var cmd int
	// var status int
	var handlers []*modbus.RTUClientHandler
	var clients []modbus.Client
	var retErrs []error
	var mqClient mqtt.Client
	_ = clients
	var err, mqErr error

	var Topics []models.Topics
	_ = Topics
	// status = 0
	ModRWChan := make(chan bool, 1)
	_ = ModRWChan
	ModChance := make(chan byte, 1)
	var wg2 sync.WaitGroup

	mqClient, mqErr = MqConnect(db)
	if mqErr != nil {
		fmt.Println("Error in Mqtt Start ", mqErr)
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Error in Mqtt Start - %s", mqErr.Error(), time.Now()))
		return handlers, clients, mqClient, &wg2, ModChance, mqErr, nil
	} else {
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Mqtt Connected", time.Now()))
		handlers, retErrs = MultiModStart(db)
		for i, err := range retErrs {
			if err != nil {
				fmt.Println("Error in Modstart ", err)
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Error in Modstart - %s,  Serial Port - %s ", time.Now(), err.Error(), handlers[i].Address))
				if i < cap(handlers) {
					handlers[i].Close()
				}

			} else {
				clients = append(clients, modbus.NewClient(handlers[i]))
				fmt.Println("Modbus Connected ")
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Modbus Connected , Serial Port - %s", time.Now(), handlers[i].Address))
			}

		}
		if cap(clients) > 0 {
			Topics, err = GetAllTopics(db)
			if err == nil {
				if len(Topics) > 0 {

					wg2.Add(len(Topics))
					// go MultiModReadWrite(mqClient, Topics, handlers, clients, ModChance, &wg2)
				} else {
					fmt.Println("There is no topics, Please add topics")
					go WsStatusPub(WsClients, fmt.Sprintf("%v	There is no topics, Please add topics", time.Now()))
				}
			} else {
				fmt.Println("Error in Fetching topics ", err)
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Error in Fetching topics ", time.Now(), err.Error()))
			}
		}
		return handlers, clients, mqClient, &wg2, ModChance, mqErr, retErrs
	}
}

func MultiModMqProcessStop(handlers []*modbus.RTUClientHandler, mqClient mqtt.Client, wg2 *sync.WaitGroup, ModChance chan byte) error {
	processMu.Lock()
	defer processMu.Unlock()
	go WsStatusPub(WsClients, fmt.Sprintf("%v	No of Gorotuines - %d", time.Now(), runtime.NumGoroutine()))
	fmt.Println("Modbus closing")
	if len(ModChance) == 0 {
		ModChance <- 0
		fmt.Println("Modbus closed")
	}
	fmt.Println("Waiting for Wait groups")
	wg2.Wait()
	fmt.Println("Finsished for Wait groups")
	for _, handler := range handlers {
		handler.Close()
	}

	mqClient.Disconnect(250)
	fmt.Println("ModChance Started")
	if len(ModChance) > 0 {
		<-ModChance
	}

	fmt.Println("Modbus Completed")
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
				topics[i].ModRegs = modregs
			}
		}
		return topics, nil
	}
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
	if err := db.Find(serialParams).Error; err != nil {
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

func ModReadWrite(mqClient mqtt.Client, Topics []models.Topics, handler *modbus.RTUClientHandler, client modbus.Client, ModChance chan byte, wg2 *sync.WaitGroup) {
	// fmt.Println("no of gorutines3 - ", runtime.NumGoroutine())
	go WsStatusPub(WsClients, fmt.Sprintf("%v	No of Gorotuines - %d", time.Now(), runtime.NumGoroutine()))
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	fmt.Println("Modbus RW Started")
	go WsStatusPub(WsClients, fmt.Sprintf("%v	Modbus RW Started", time.Now()))

	for _, topic := range Topics {
		payload := make(chan string)
		mqShutDown := make(chan bool)
		go MqttOperation(mqClient, topic, payload, mqShutDown, wg2)
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Goroutine STARTED for Publishing Mqtt Topic - %s", time.Now(), topic.Topic))
		go ModOperation(topic, handler, client, ModChance, wg2, payload, mqShutDown)
		go WsStatusPub(WsClients, fmt.Sprintf("%v	Goroutine STARTED for Reading Modbus Registers - %v which are related to Topic - %s", time.Now(), topic.ModRegs, topic.Topic))
	}
	ModChance <- 100
	return
}

func ModOperation(topic models.Topics, handler *modbus.RTUClientHandler, client modbus.Client, ModChance chan byte, wg2 *sync.WaitGroup, payload chan string, mqShutDown chan bool) {
	// fmt.Println("Modbus connected and commanded to fetch data")
	for {
		select {
		case chance := <-ModChance:
			if chance == 100 {
				for _, modreg := range topic.ModRegs {
					time.Sleep(time.Duration(topic.Delay) * time.Second)
					handler.SlaveId = modreg.Unit

					results, err := client.ReadHoldingRegisters(modreg.Register, modreg.Qty)

					if err == nil {
						if modreg.DataType == 1 {
							// uint 8
							val := results
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							// payload <- fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %f, err - %v", topic.Topic, modreg.Name, modreg.Register, val, err)
							if err == nil {

								payload <- string(fmt.Sprintf("{\"%s\" : \"%d\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %d`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
							}
						} else if modreg.DataType == 2 {
							// []uint 8
						} else if modreg.DataType == 3 {
							// int 8
						} else if modreg.DataType == 4 {
							// []int 8
						} else if modreg.DataType == 5 {
							// uint16
							val, err := SingleUint16FromBytes(results, modreg.ByteOrder)
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							if err == nil {
								// payload <- fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val)
								payload <- string(fmt.Sprintf("{\"%s\" : \"%d\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %f`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
							}

						} else if modreg.DataType == 6 {
							// []uint16
							val, err := arrUint16frombytes(results, modreg.ByteOrder)
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							if err == nil {
								// payload <- fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val)
								payload <- string(fmt.Sprintf("{\"%s\" : \"%d\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %f`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
							}

						} else if modreg.DataType == 7 {
							// int16
							val, err := SingleUint32FromBytes(results, modreg.ByteOrder)
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							if err == nil {
								// payload <- fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val)
								payload <- string(fmt.Sprintf("{\"%s\" : \"%d\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %f`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
							}

						} else if modreg.DataType == 8 {
							// []int16
							val, err := arrUint32frombytes(results, modreg.ByteOrder)
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %d, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							if err == nil {
								// payload <- fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val)
								payload <- string(fmt.Sprintf("{\"%s\" : \"%d\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %f`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
							}
						} else if modreg.DataType == 9 {
							// uint32
						} else if modreg.DataType == 10 {
							// []uint32
						} else if modreg.DataType == 11 {
							// int32
						} else if modreg.DataType == 12 {
							// []int32
						} else if modreg.DataType == 13 {
							// uint64
						} else if modreg.DataType == 14 {
							// []uint64
						} else if modreg.DataType == 15 {
							// int64
						} else if modreg.DataType == 16 {
							// []int64
						} else if modreg.DataType == 17 {
							// float32
							val, err := SingleFloat32FromBytes(results, modreg.ByteOrder)
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %f, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							if err == nil {
								// payload <- fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val)
								payload <- string(fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %f`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
							}

						} else if modreg.DataType == 18 {
							// []float32
							val, err := arrFloat32frombytes(results, modreg.ByteOrder)
							go WsClientPub(ModLastAquired{fmt.Sprintf("Topic- %s, Reg Name- %s, Register - %d, Value %f, err- %v ", topic.Topic, modreg.Name, modreg.Register, val, err)})
							if err == nil {
								// payload <- fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val)
								payload <- string(fmt.Sprintf("{\"%s\" : \"%f\"}", modreg.Tags, val))
								// content, b := json.Marshal(JSONString(fmt.Sprintf(`"%s" : %f`, modreg.Tags, val)))
								// if b != nil {
								// 	payload <- content
								// }
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
						fmt.Println(topic.Topic, modreg.Name, modreg.Register, results, err)
					}
				}
				// Concat payload and Send as Combined Payload here
				if len(ModChance) == 0 {

					ModChance <- 100
				}
			} else {
				if len(ModChance) == 0 {
					ModChance <- 255
					mqShutDown <- true
				}
				// wg2.Done()
				fmt.Println("Topic - ", topic.Topic, " goroutine is stopped")
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Goroutine STOPPED to Reading Modbus Registers - %v which are related to Topic - %s", time.Now(), topic.ModRegs, topic.Topic))

				return

			}

		}

	}

}

func MqttOperation(mqClient mqtt.Client, topic models.Topics, payload chan string, mqShutDown chan bool, wg2 *sync.WaitGroup) {
	for {
		select {
		case pay := <-payload:

			mqClient.Publish(topic.Topic, 0, false, pay)
			fmt.Sprintf("Topic - %s, Payload- %s", topic.Topic, pay)
			go WsClientPub(MqLastSent{fmt.Sprintf("Topic - %s, Payload- %s", topic.Topic, pay)})
		case cmd := <-mqShutDown:
			if cmd {
				fmt.Println("Shutdown mq operation  with topic -", topic.Topic)
				go WsStatusPub(WsClients, fmt.Sprintf("%v	Goroutine STOPPED to Publishing Mqtt Topic - %s", time.Now(), topic.Topic))
				wg2.Done()
				return
			}

		}
	}
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

	return gval.Evaluate(ppstring, map[string]interface{}{"value": value})

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
