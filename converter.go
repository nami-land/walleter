package walleter

import (
	"fmt"
	"strconv"
	"strings"
)

type Number interface {
	int64 | uint64 | float64
}

func convertArrayToString[T Number](array []T, symbol string) string {
	var temp = make([]string, len(array))
	for k, v := range array {
		temp[k] = fmt.Sprintf("%d", v)
	}
	return strings.Join(temp, symbol)
}

// indexOfArray get index of value in array. if value doesn't exist, return -1.
func indexOfArray[T Number](array []T, value T) int {
	for index, num := range array {
		if num == value {
			return index
		}
	}
	return -1
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
