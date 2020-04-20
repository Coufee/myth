package utils

import (
	"bytes"
	"reflect"
	"strings"
)

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func SplitSubN(s string, n int) []string {
	sub := ""
	var subs []string

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}

//传入通道时只检查是否有数据
func CheckIn(search interface{}, origin interface{}) bool {
	switch reflect.TypeOf(origin).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(origin)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == search {
				return true
			}
		}
	case reflect.Map:
		s := reflect.ValueOf(origin)
		for _, value := range s.MapKeys() {
			if s.MapIndex(value).Interface() == search {
				return true
			}
		}
	case reflect.Chan:
		s := reflect.ValueOf(origin)
		if s.Len() > 0 {
			return true
		}
	case reflect.String:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Bool:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Int:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Int8:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Int16:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Int32:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Int64:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Float32:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	case reflect.Float64:
		s := reflect.ValueOf(origin)
		if s.Interface() == search {
			return true
		}
	}

	return false
}

func CheckInArray(search interface{}, list []interface{}) bool {
	for _, v := range list {
		if search == v {
			return true
		}
	}

	return false
}

//func CheckStringInSlice(search string, list []string) bool {
//	for _, key := range list {
//		if key == search {
//			return true
//		}
//	}
//
//	return false
//}
//
//func CheckIntInSlice(search int, list []int) bool {
//	for _, key := range list {
//		if key == search {
//			return true
//		}
//	}
//
//	return false
//}
//
//func CheckInt32InSlice(search int32, list []int32) bool {
//	for _, key := range list {
//		if key == search {
//			return true
//		}
//	}
//
//	return false
//}
//
//func CheckInt64InSlice(search int64, list []int64) bool {
//	for _, key := range list {
//		if key == search {
//			return true
//		}
//	}
//
//	return false
//}

func IntSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func Int32SliceEqual(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func Int64SliceEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func Float32SliceEqual(a, b []float32, precision int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if precision > 0 {
			if DecimalFloat32(v, precision) != DecimalFloat32(b[i], precision) {
				return false
			}
		} else {
			if v != b[i] {
				return false
			}
		}
	}

	return true
}

func Float64SliceEqual(a, b []float64, precision int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if precision > 0 {
			if DecimalFloat64(v, precision) != DecimalFloat64(b[i], precision) {
				return false
			}
		} else {
			if v != b[i] {
				return false
			}
		}
	}

	return true
}

func Int32SliceToStringSlice(origin []int32) []string {
	var list = make([]string, 0)

	for _, item := range origin {
		list = append(list, NumToString(item))
	}

	return list
}

func StringToIntSlice(origin, separator string) ([]int, error) {
	intSlice := make([]int, 0)
	if origin != "" {
		for _, intStr := range strings.Split(origin, separator) {
			num, err := StringToInt(intStr)
			if err != nil {
				return intSlice, err
			} else {
				intSlice = append(intSlice, num)
			}
		}
	}

	return intSlice, nil
}

func StringToInt32Slice(origin, separator string) ([]int32, error) {
	int32Slice := make([]int32, 0)
	if origin != "" {
		for _, intStr := range strings.Split(origin, separator) {
			num, err := StringToInt32(intStr)
			if err != nil {
				return int32Slice, err
			} else {
				int32Slice = append(int32Slice, num)
			}
		}
	}

	return int32Slice, nil
}

func StringToInt64Slice(origin, separator string) ([]int64, error) {
	int64Slice := make([]int64, 0)
	if origin != "" {
		for _, intStr := range strings.Split(origin, separator) {
			num, err := StringToInt64(intStr)
			if err != nil {
				return int64Slice, err
			} else {
				int64Slice = append(int64Slice, num)
			}
		}
	}

	return int64Slice, nil
}

func ArrayIntUnique(l []int) []int {
	m := make(map[int]int)
	for _, v := range l {
		m[v] = v
	}

	r := make([]int, 0)
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

func ArrayInt32Unique(l []int32) []int32 {
	m := make(map[int32]int32)
	for _, v := range l {
		m[v] = v
	}

	r := make([]int32, 0)
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

func ArrayInt64Unique(l []int64) []int64 {
	m := make(map[int64]int64)
	for _, v := range l {
		m[v] = v
	}

	r := make([]int64, 0)
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

func ArrayStringUnique(l []string) []string {
	m := make(map[string]string)
	for _, v := range l {
		m[v] = v
	}

	r := make([]string, 0)
	for _, v := range m {
		r = append(r, v)
	}

	return r
}
