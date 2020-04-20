package utils

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

func BytesToStringFast(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//占位补全
func StringComplete(origin, placeholder string, maxLen int, prefix bool) string {
	curLen := len(origin)
	if curLen >= maxLen {
		return origin[0:maxLen]
	} else {
		if prefix {
			temp := strings.Repeat(placeholder, maxLen-curLen)
			return temp + origin
		} else {
			temp := strings.Repeat(placeholder, maxLen-curLen)
			return origin + temp
		}
	}

	return ""
}

//反转
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//整数转string
func NumToString(num interface{}) string {
	typeStr := reflect.TypeOf(num)
	switch typeStr.Kind() {
	case reflect.Int:
		return strconv.FormatInt(int64(num.(int)), 10)
	case reflect.Int8:
		return strconv.FormatInt(int64(num.(int8)), 10)
	case reflect.Int16:
		return strconv.FormatInt(int64(num.(int16)), 10)
	case reflect.Int32:
		return strconv.FormatInt(int64(num.(int32)), 10)
	case reflect.Int64:
		return strconv.FormatInt(num.(int64), 10)
	case reflect.String:
		return num.(string)
	default:
		return ""
	}

	return ""
}

var (
	bfPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer([]byte{})
		},
	}
)

func Int32SliceToString(origin []int32, separator string) string {
	var list = make([]string, 0)
	for _, item := range origin {
		list = append(list, NumToString(item))
	}

	return JoinStrings(separator, list...)
}

func MergeStrings(stringArray ...string) string {
	var buffer bytes.Buffer
	for _, v := range stringArray {
		buffer.WriteString(v)
	}
	return buffer.String()
}

func JoinStrings(separator string, stringArray ...string) string {
	//var buffer bytes.Buffer
	buf := bfPool.Get().(*bytes.Buffer)
	var max = len(stringArray) - 1
	for vi, v := range stringArray {
		buf.WriteString(v)
		if vi < max {
			buf.WriteString(separator)
		}
	}

	s := buf.String()
	buf.Reset()
	bfPool.Put(buf)
	return s
}
