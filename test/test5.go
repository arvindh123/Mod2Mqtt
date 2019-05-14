package main

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/PaesslerAG/gval"
)

type myFloat64 float64
type KeepZero float64

func (mf myFloat64) MarshalJSON() ([]byte, error) {
	const ε = 1e-12
	v := float64(mf)
	w, f := math.Modf(v)
	if f < ε {
		return []byte(fmt.Sprintf(`%v.0`, math.Trunc(w))), nil
	}
	return json.Marshal(v)
}

func (f KeepZero) MarshalJSON() ([]byte, error) {
    if float64(f) == float64(int(f)) {
        return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
    }
    return []byte(strconv.FormatFloat(float64(f), 'f', -1, 32)), nil
}

func main() {

	// value, err := gval.Evaluate("(mem_used / total_mem) ",
	// 	map[string]interface{}{
	// 		"total_mem": 1024.234,
	// 		"mem_used":  512.4234,
	// 	})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Print(value)
	var a map[string]interface{} 

	// fmt.Println(a[2])
	// ArrPostProcess("value / 10 :: value / 100 :: value / 1000", a)
	val, err := SinglePostProcess("(value/100)", 500)
	a["h"] = val
	switch 	.(type) {
	case float64:
		fmt.Println("float64")
		a["h"] =
	case string:
		fmt.Println("string")
	default:
		fmt.Println("others")
	}
	fmt.Println(val, err)
}

func SinglePostProcess(ppstring string, value interface{}) (interface{}, error) {

	value, err := gval.Evaluate(ppstring, map[string]interface{}{"value": value})
	switch value.(type) {
	case float64:
		value = myFloat64{fmt.Sprintf(`%.2f`,value)
		// const ε = 1e-12
		// fmt.Println(`ε -`, ε)
		// fmt.Println("float64")
		// v := value.(float64)
		// w, f := math.Modf(v)
		// if f < ε {
		// 	fmt.Println(`f < ε`)
		// 	fmt.Println([]byte(fmt.Sprintf(`%.2f`, math.Trunc(w))), nil)

		// } else {
		// 	fmt.Println(`else`)
		// 	fmt.Println(w, f)
		// }
	default:
		fmt.Println("others")
	}
	fmt.Println(value, err)
	return value, err
}
