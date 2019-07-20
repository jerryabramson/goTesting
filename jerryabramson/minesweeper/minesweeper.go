/*
 * Simple implementation of MineSwepper:
 *     https://en.wikipedia.org/wiki/Minesweeper_(video_game)
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	minesArg := int(9)
	widthArg := int(9)
	heightArg := int(9)
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
			widthArg, err = SafeAtoI(dim[0])
			heightArg, err = SafeAtoI(dim[1])
		} else if strings.Compare(arg, "--mines") == 0 {
			if strings.Compare(argVal, "NONE") == 0 {
				Usage(nil)
			}
			minesArg, err = SafeAtoI(argVal)
			argc++
		} else {
			Usage(nil)
		}
		argc++
	}
	b := New(minesArg, widthArg, heightArg)
	PopulateBoard(&b)
	Play(&b)
}

func Usage(err error) {
	if err == nil {
		err = errors.New("")
	}
	fmt.Printf("\nUsage: MineSweeper [--size <x,y>] --mines <numberOfMines>: %v\n", err)
	os.Exit(1)
}
