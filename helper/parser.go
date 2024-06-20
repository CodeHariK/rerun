package helper

import (
	"strconv"
	"strings"
)

// Convert string "1,2,3" to int array []int{1,2,3}
func ParseStringInts(value string) ([]int, error) {
	var i []int
	for _, v := range strings.Split(value, ",") {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		i = append(i, intValue)
	}
	return i, nil
}
