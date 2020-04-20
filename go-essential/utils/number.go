package utils

import (
	"fmt"
	"strconv"
)

//省略小数点后几位
func DecimalFloat32(value float32, precision int) float32 {
	result, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(precision)+"f", value), 32)
	return float32(result)
}

//省略小数点后几位
func DecimalFloat64(value float64, precision int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(precision)+"f", value), 64)
	return value
}

func StringToInt(str string) (int, error) {
	intParam, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return intParam, nil
}

func StringToInt32(str string) (int32, error) {
	int32Param, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(int32Param), nil
}

func StringToInt64(str string) (int64, error) {
	int64Param, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return int64Param, nil
}
