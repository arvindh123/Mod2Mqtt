package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PaesslerAG/gval"
)

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
	a := 5
	// fmt.Println(a[2])
	ArrPostProcess("value * 2", a)
}

func ArrPostProcess(ppstring string, values ...interface{}) {
	ppstrings := strings.Split(ppstring, "::")
	capPpstrings := cap(ppstrings)

	for _, value := range values {
		switch reflect.TypeOf(value).Kind() {

		case reflect.Slice:
			s := reflect.ValueOf(value)
			fmt.Println("Slice\n")
			fmt.Println(s)
			for i := 0; i < s.Len(); i++ {
				if i < capPpstrings {
					// fmt.Println(s.Index(i) * 2)
					SinglePostProcess(ppstrings[i], s.Index(i))
				} else {
					SinglePostProcess(ppstrings[0], s.Index(i))
				}
			}

		case reflect.Array:
			s := reflect.ValueOf(value)
			fmt.Println("Slice\n")
			fmt.Println(s)

		default:
			s := reflect.ValueOf(value)
			fmt.Println(s)
			SinglePostProcess(ppstrings[0], value)
		}

	}

}

func SinglePostProcess(ppstring string, value interface{}) (interface{}, error) {
	value, err := gval.Evaluate(ppstring, map[string]interface{}{"value": value})
	fmt.Println(value, err)
	return value, err

}
