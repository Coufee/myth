package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func CreateNumberRandCode(len int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < len; i++ {
		num := rand.Intn(9)
		code = code + strconv.Itoa(num)
	}
	return code
}

var (
	lowerSlice   = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	capitalSlice = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func CreateCapitalLetterRandCode(count int) string {
	var axis string
	for count > 0 {
		originSize := len(capitalSlice)
		if count%originSize == 0 {
			axis = capitalSlice[originSize-1] + axis
			count -= originSize
		} else {
			axis = capitalSlice[(count%originSize)-1] + axis
		}
		count /= originSize
	}

	return axis
}

func CreateLowerLetterRandCode(count int) string {
	var axis string
	for count > 0 {
		originSize := len(lowerSlice)
		if count%originSize == 0 {
			axis = lowerSlice[originSize-1] + axis
			count -= originSize
		} else {
			axis = lowerSlice[(count%originSize)-1] + axis
		}
		count /= originSize
	}

	return axis
}

func CreateLetterRandCode(count int) string {
	var axis string
	for count > 0 {
		originSlice := make([]string, 0)
		originSlice = append(originSlice, capitalSlice...)
		originSlice = append(originSlice, lowerSlice...)
		originSize := len(originSlice)
		if count%originSize == 0 {
			axis = lowerSlice[originSize-1] + axis
			count -= originSize
		} else {
			axis = lowerSlice[(count%originSize)-1] + axis
		}
		count /= originSize
	}

	return axis
}
