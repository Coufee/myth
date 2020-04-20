package utils

import (
	"fmt"
	"testing"
)

func TestDecimal(t *testing.T) {
	data := 1.234
	result := DecimalFloat64(data,2)
	fmt.Print(result)
}

func TestCheckInArray(t *testing.T) {
	//data := []interface{}{"a","b","c"}
	//data := "b"
	//data := map[int]string{1:"a",2:"b1"}
	//data := make(chan string,10)
	//data<-"a"
	data := 32
	result := CheckIn(1, data)
	fmt.Printf("%v \n", result)
}

