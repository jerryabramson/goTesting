package main

import (
	"strconv"
	"errors"
	"strings"
	"fmt"
)

var (
	    err error   
)

func SafeAtoI(val string) (int, error) {
	ret, e := strconv.Atoi(val)
	if e == nil {
		return ret, e
	}
	var errMsg strings.Builder
	fmt.Fprintf(&errMsg, "utils: Invalid Number '%s'", val)
	e = errors.New(errMsg.String())
	return ret, e
}

func Err() error {
	return err
}
