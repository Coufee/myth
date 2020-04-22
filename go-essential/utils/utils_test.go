package utils

import (
	"testing"
	"fmt"
)

//func TestDecimal(t *testing.T) {
//	data := 1.234
//	result := DecimalFloat64(data,2)
//	fmt.Print(result)
//}
//
//func TestCheckInArray(t *testing.T) {
//	//data := []interface{}{"a","b","c"}
//	//data := "b"
//	//data := map[int]string{1:"a",2:"b1"}
//	//data := make(chan string,10)
//	//data<-"a"
//	data := 32
//	result := CheckIn(1, data)
//	fmt.Printf("%v \n", result)
//}

type Person interface {
	Name() string
}

type ChenQiongHe struct {
}

func (t *ChenQiongHe) Name() string {
	return "雪山飞猪"
}

func TestVerifyNil(t *testing.T) {
	//var a interface{}
	//result := VerifyNil(a)
	//fmt.Print("------------------",result)
	//a=nil
	//result = VerifyNil(a)
	//fmt.Print("------------------",result)
	var test *ChenQiongHe
	if test == nil {
		fmt.Println("test == nil")
	} else {
		fmt.Println("test != nil")
	}
	//将空指针赋值给接口
	var person Person = test
	if VerifyNil(person) {
		fmt.Print("person == nil")
	} else {
		fmt.Print("person != nil")
	}
}

