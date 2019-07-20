package main

import (
	"errors"
	"strconv"
)

func SafeAtoI(val string) int {
	ret := int(-1)
	err := errors.New("none")
	ret, err = strconv.Atoi(val)
	if (err != nil) {
		ret = -1
	}
	return ret
}
