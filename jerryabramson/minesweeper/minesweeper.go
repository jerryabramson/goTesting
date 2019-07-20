package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/jerryabramson/minesweeper/libs"
)

// Main entry point
func main() {
	minesArg := int(9)
	widthArg := int(9)
	heightArg := int(9)
	var err error
	for argc := 1; argc < len(os.Args); {
		arg := os.Args[argc]
		argVal := "NONE"
		if argc < len(os.Args)-1 {
			argVal = os.Args[argc+1]
		}
		if strings.Compare(arg, "--size") == 0 {
			if strings.Compare(argVal, "NONE") == 0 {
				Usage(nil)
			}
			dim := strings.Split(argVal, ",")
			if len(dim) != 2 {
				Usage(nil)
			}
			argc++
			widthArg, err = minesweeper.SafeAtoI(dim[0])
			if (err != nil) {
				Usage(err)
			}
			heightArg, err = minesweeper.SafeAtoI(dim[1])
			if (err != nil) {
				Usage(err)
			}
		} else if strings.Compare(arg, "--mines") == 0 {
			if strings.Compare(argVal, "NONE") == 0 {
				Usage(nil)
			}
			minesArg, err = minesweeper.SafeAtoI(argVal)
			if (err != nil) {
				Usage(err)
			}
			argc++
		} else {
			Usage(nil)
		}
		argc++
	}
	b := minesweeper.New(minesArg, widthArg, heightArg)
	minesweeper.PopulateBoard(&b)
	minesweeper.Play(&b)
}

func Usage(err error) {
	if err == nil {
		err = errors.New("")
	}
	fmt.Printf("\nUsage: MineSweeper [--size <x,y>] --mines <numberOfMines>: %v\n", err)
	os.Exit(1)
}
