package wallet_center

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func convertUintArrayToString(arr []uint64, symbol string) string {
	var temp = make([]string, len(arr))
	for k, v := range arr {
		temp[k] = fmt.Sprintf("%d", v)
	}
	return strings.Join(temp, symbol)
}

func convertStringToUIntArray(str string) []uint64 {
	var result []uint64
	strArray := strings.Split(str, ",")
	for _, item := range strArray {
		intItem, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			continue
		}
		result = append(result, uint64(intItem))
	}
	return result
}

func convertIntArrayToBigIntArray(arr []int32) []*big.Int {
	var result []*big.Int
	for _, item := range arr {
		result = append(result, big.NewInt(int64(item)))
	}
	return result
}

func convertUInt64ArrayToUIntArray(arr []uint64) []uint {
	var result []uint
	for _, item := range arr {
		result = append(result, uint(item))
	}
	return result
}

func getIndexFromUIntArray(array []uint64, value uint64) int {
	for index, num := range array {
		if num == value {
			return index
		}
	}
	return -1
}
