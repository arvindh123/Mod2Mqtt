package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

type MyF64 float64

// func (f MyF64) MarshalJSON() ([]byte, error) {
// 	defer elapsed("Marshal MyF64 to byte - ")()
// 	if float64(f) == float64(int(f)) {
// 		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
// 	}
// 	return json.Marshal(f)
// }

func (mf MyF64) MarshalJSON() ([]byte, error) {
	const ε = 1e-12
	v := float64(mf)
	w, f := math.Modf(v)
	if f < ε {
		return []byte(fmt.Sprintf(`%v.0`, math.Trunc(w))), nil
	}
	return json.Marshal(v)
}

func main() {
	var a map[string]interface{}
	a = make(map[string]interface{})
	a["h"] = 3.0
	fmt.Println("Start time", time.Now())
	switch a["h"].(type) {
	case float64:

		fmt.Println("float64", a["h"])
		fmt.Println("Start of MyF64 conversion time", time.Now())
		a["h"] = MyF64(a["h"].(float64))
		fmt.Print("End of MyF64 conversion time", time.Now())
	case MyF64:
		fmt.Println("MyF64 ", a["h"])
	}

	switch a["h"].(type) {
	case float64:
		fmt.Println("float64", a["h"])
		//a["h"] =  MyF64(a["h"].(float64))
	case MyF64:
		fmt.Println("MyF64 ", a["h"])
	}
	a["w"] = 324
	fmt.Println("Start of Marshal conversion time", time.Now())
	fmt.Println(json.Marshal(a))
	fmt.Println("End of Marshal conversion time", time.Now())
	fmt.Println("End time", time.Now())
	fmt.Println("Hello, playground")
}
