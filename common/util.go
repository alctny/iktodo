package common

import (
	"strconv"
)

// StringsToInts Conver []string to []int
func StringsToInts(args []string) ([]int, error) {
	res := []int{}
	for _, arg := range args {
		i, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, int(i))
	}
	return res, nil
}
