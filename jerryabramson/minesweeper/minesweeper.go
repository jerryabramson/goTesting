/*
 * Simple implementation of MineSwepper:
 *     https://en.wikipedia.org/wiki/Minesweeper_(video_game)
 */

// leetcode problem to find median of two integer arrays
package main

import (
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
		if (argc < len(os.Args) - 1) {
			argVal = os.Args[argc+1];
		}
//		fmt.Printf("argc = %d, arg = %s, argVal = %s\n",
//			argc, arg, argVal);
		if (strings.Compare(arg, "--size") == 0) {
			if (strings.Compare(argVal, "NONE") == 0) {
				Usage()
			}
			dim := strings.Split(argVal, ",")
			if (len(dim) != 2) {
				Usage();
			}
			argc++
			widthArg = SafeAtoI(dim[0]);
			heightArg = SafeAtoI(dim[1]);
		} else if (strings.Compare(arg, "--mines") == 0) {
			if (strings.Compare(argVal, "NONE") == 0) {
				Usage();
			}
			minesArg = SafeAtoI(argVal);
			argc++
		} else {
			Usage();
		}
		argc++
	}
	b := New(minesArg, widthArg, heightArg)
	PopulateBoard(&b)
	Play(&b)
}

func Usage() {
	fmt.Println("Usage: MineSweeper [--size <x,y>] --mines <numberOfMines>");
	os.Exit(1);
}
    
