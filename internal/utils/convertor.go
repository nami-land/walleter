package utils

import (
	"math/big"
	"strconv"
	"strings"
)

func ConvertUintArrayToString(arr []uint) string {
	var result = ""
	for i := 0; i < len(arr); i++ {
		if i == len(arr)-1 {
			result += strconv.Itoa(int(arr[i]))
		} else {
			result += strconv.Itoa(int(arr[i])) + ","
		}
	}
	return result
}

func ConvertStringToIntArray(str string) []int32 {
	var result []int32
	strArray := strings.Split(str, ",")
	for _, item := range strArray {
		intItem, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			continue
		}
		result = append(result, int32(intItem))
	}
	return result
}

func ConvertIntArrayToBigIntArray(arr []int32) []*big.Int {
	var result []*big.Int
	for _, item := range arr {
		result = append(result, big.NewInt(int64(item)))
	}
	return result
}
