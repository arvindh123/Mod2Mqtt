package mod

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/arvindh123/Mod2Mqtt/ws"
)

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
	switch val.(type) {
	case float64:
		if val == math.Trunc(val.(float64)) {
			val = myF64(val.(float64))
		}
	}
	if err != nil {
		go ws.WsStatusPub(fmt.Sprintf("Error in Post Process Ppstring -%s Value -%v", ppstring, value))
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
