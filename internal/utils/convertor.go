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

func ConvertStringToUIntArray(str string) []uint {
	var result []uint
	strArray := strings.Split(str, ",")
	for _, item := range strArray {
		intItem, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			continue
		}
		result = append(result, uint(intItem))
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

func ConvertUInt64ArrayToUIntArray(arr []uint64) []uint {
	var result []uint
	for _, item := range arr {
		result = append(result, uint(item))
	}
	return result
}

func CovertUIntArrayToUInt64Array(arr []uint) []uint64 {
	var result []uint64
	for _, item := range arr {
		result = append(result, uint64(item))
	}
	return result
}

func GetIndexFromUIntArray(array []uint, value uint) int {
	for index, num := range array {
		if num == value {
			return index
		}
	}
	return -1
}
