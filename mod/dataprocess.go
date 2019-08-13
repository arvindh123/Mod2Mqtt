package mod

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/arvindh123/Mod2Mqtt/ws"
)

type myF64 float64

func (f myF64) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
	}
	return json.Marshal(f)
}

func ModReadDataProcess(Device models.DeviceDetails, modreg models.ModbusRegisters, results []byte, Payload chan map[string]interface{}, ModWg *sync.WaitGroup, err error) {
	retString := make(map[string]interface{})
	if err == nil {
		if modreg.DataType == 1 {
			// 1 - uint 8 // 3 - int 8
			val := results
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
		} else if modreg.DataType == 2 {
			// []uint 8
			val := results
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			} else {
				val, err := SingleUint16FromBytes(results, modreg.ByteOrder)
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			}
		} else if modreg.DataType == 6 {
			// []uint16
			if strings.Contains(modreg.PostProcess, "value") {
				val, err := arrUint16frombytesPP(results, modreg.ByteOrder, modreg.PostProcess)
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			} else {
				val, err := SingleUint16FromBytes(results, modreg.ByteOrder)
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
				if err == nil {
					retString[modreg.Tags] = val
				}
			}
		} else if modreg.DataType == 8 {
			// []int16
			if strings.Contains(modreg.PostProcess, "value") {
				val, err := arrUint16frombytesPP(results, modreg.ByteOrder, modreg.PostProcess)
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
				go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}
		} else if modreg.DataType == 12 {
			// []int32
			val, err := arrUint32frombytes(results, modreg.ByteOrder)
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}
		} else if modreg.DataType == 16 {
			// []int64
			val, err := arrUint64frombytes(results, modreg.ByteOrder)
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %f, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
			if err == nil {
				retString[modreg.Tags] = val
			}

		} else if modreg.DataType == 18 {
			// []float32
			val, err := arrFloat32frombytes(results, modreg.ByteOrder)
			go ws.WsClientPub(ws.ModLastAquired{fmt.Sprintf("Device- %s, Reg Name- %s, Register - %d, Value %f, err- %v ", Device.Name, modreg.Name, modreg.Register, val, err)})
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
		Payload <- retString
		ModWg.Done()
	} else {
		ModWg.Done()
		go ws.WsStatusPub(fmt.Sprintf("%v	Device- %s, Reg Name- %s, Register - %d, Value %v, err- %v ", time.Now(), Device.Name, modreg.Name, modreg.Register, results, err))
		// fmt.Println(Device.Name, modreg.Name, modreg.Register, results, err)
	}
	return
}
