package minesweeper

import (
	"fmt"
	"strconv"
	"time"
	"bufio"
	"io"
	"strings"
	"os"
)

const (
    ANSI_LINE                   = "\u2500"
    BOTTOM_LEFT_CORNER          = "\u2514"
    BOTTOM_RIGHT_CORNER         = "\u2518"
    LEFT_COLUMN_LINE            = "\u251c"
    BOTTOM_VERTICAL_BAR         = "\u253C"
    TOP_VERTICAL_BAR            = "\u252C"
    BOTTOM_BAR                  = "\u2534"
    RIGHT_COLUMN_LINE           = "\u2524"
    TOP_LEFT_CORNER             = "\u250c"
    TOP_RIGHT_CORNER            = "\u2510"
    VERTICAL_BAR                = "\u2502"
	ORIGINY                     = 4
	ORIGINX                     = 8
)

// Reveal either the board discovered so far,
// or the underlying board itself (after win/lose)
func drawWholeBoard(b *board, h bool) {
	fmt.Printf("\033[0H\033[2J")
	fmt.Printf("\033[0H\033[1mCurrent Board\n\033[0m")
	printLine(80)
	fmt.Println()
	fmt.Printf("\033[0m      ")
	for x := 0; x < b.width; x++ {
		fmt.Printf("%-2d  ", x+1)
	}
	fmt.Printf("\n")
	fmt.Printf("    %s", TOP_LEFT_CORNER) 
	for x := 0; x < b.width - 1; x++ {
		printLine(3)
		fmt.Printf(TOP_VERTICAL_BAR)
	}
	printLine(3)
	fmt.Printf(TOP_RIGHT_CORNER)
	fmt.Printf("\n")
	for y := 0; y < b.height; y++ {
		fmt.Printf("%4d%s", y+1, VERTICAL_BAR)
		for x := 0; x < b.width; x++ {
			fmt.Printf("%s%s", showPiece(b, x, y, h), VERTICAL_BAR)
		}
		fmt.Println()
		if (y < b.height - 1) {
			fmt.Printf("    %s", LEFT_COLUMN_LINE) 
			for x := 0; x < b.width - 1; x++ {
				printLine(3)
				fmt.Printf(BOTTOM_VERTICAL_BAR)
			}
			printLine(3)
			fmt.Printf(RIGHT_COLUMN_LINE)
			fmt.Printf("\n")
		}
		
	}
	fmt.Printf("    %s", BOTTOM_LEFT_CORNER) 
	for x := 0; x < b.width - 1; x++ {
		printLine(3)
		fmt.Printf(BOTTOM_BAR)
	}
	printLine(3)
	fmt.Printf(BOTTOM_RIGHT_CORNER)
	fmt.Println()
	fmt.Printf("\033[0m     ")
	for x := 0; x < b.width; x++ {
		fmt.Printf("%2d  ", x+1)
	}
	fmt.Println()
}

// fancy display of a point on the board.
func showPiece(b *board, x int, y int, h bool) string {
	myPieces := b.discoveredBoard
	if !h {
		myPieces = b.actualBoard
	}
	p := myPieces[y][x]
	if p == MINE {
		return "\033[41;1m \033[0m\033[31;1mX\033[0m\033[41;1m \033[0m"
	} else if p == UNKNOWN {
		return "\033[47m - \033[0m"
	} else if p == POSSIBLE {
		return "\033[43m ? \033[0m"
	} else {
		ps := strconv.FormatInt(p, 10)
		pString := string("X")
		if p > 0 {
			pString = "\033[42m " + ps + " \033[0m"
		} else {
			pString = "\033[46m " + ps + " \033[0m"
		}
		return pString
	}
}


func revealPiece(b *board, x int, y int, h bool) {
	row := ORIGINY + y * 2 + 1
	column := ORIGINX + (x-1) * 4 + 2
	fmt.Printf("\033[%d;%dH%s", row, column, showPiece(b, x,y, h))
	clearStatus()
	showLocation(b, x, y)
}

func clearStatus() {
	fmt.Printf("\033[1;16H\033[K")
}

func showLocation(b *board, x int, y int) {
	fmt.Printf("\033[1;74H\033[47m%02d\033[0m , \033[47m%02d\033[0m", y+1, x+1)
}

func prompt(b *board) (int, int, error, bool) {
	x := -1
	y := -1
	ierr := error(nil)
	sc := bufio.NewScanner(os.Stdin)
	safeMark := bool(false)
	
	fmt.Printf("\033[24;1HPlease choose a spot to check for a mine [origin at 1,1] (y,x)[,?]:\033[K ")
	err := sc.Scan()
	if !err {
		msg := sc.Err()
		if msg == nil {
			msg = io.EOF
		}
		fmt.Printf("\033[1;16H\033[31mI/O Error: %v\033[K\033[m", msg)
		fmt.Printf("\033[26;1H")
		x=-99
		y=-99
	} else {
		input := sc.Text()
		dim := strings.Split(input, ",")
		if len(dim) == 3 {
			if strings.Compare(dim[2], "?") == 0 {
				safeMark = true;
			}
		} else if len(dim) != 2 {
			fmt.Printf("\033[1;16H\033[31mInvalid Syntax\033[K\033[0m")
			x=-1
			y=-1
		}
		y, ierr := SafeAtoI(dim[0])
		y--
		if ierr != nil {
			fmt.Print("\033[1;16H\033[31m%s\033[K\033[0m", ierr)
			x=-1
			y=-1
		} else {
			x, ierr := SafeAtoI(dim[1])
			x--
			if ierr != nil {
				fmt.Printf("\033[1;16H\033[31m%s\033[K\033[0m", ierr)
				x=-1
				y=-1
			} else {
				return x,y, ierr, safeMark
			}
		}
	}
	return x,y, ierr, safeMark
}


func resetSpace (b *board, x int, y int) {
	fmt.Printf("\033[1;16H\033[32mresetting space %d,%d from\033[0m %s ",
		x+1, y+1, showPiece(b, x, y, true))
	fmt.Printf(" \033[32mback to\033[0m %s\033[K", showPiece(b,x,y,true))
	showLocation(b, x, y)
}

func duplicateMove(b *board, x int, y int) {
	fmt.Printf("\033[1;16H\033[35mYou have already revealed space %d,%d: value = \033[0m%s\033[K",
		y+1, x+1, showPiece(b, x, y, true))
	showLocation(b, x, y)
}

func outOfRange() {
	fmt.Printf("\033[1;16H\033[31m** Out of Range\033[K\033[0m")
}



func win(b *board, x int, y int) {
	fmt.Printf("\033[0H\033[2J")
	fmt.Printf("\033[0H\033[1mCurrent Board\n\033[0m")
	printLine(80)
	fmt.Println()
	drawWholeBoard(b, false)
	fmt.Println()
	fmt.Printf("\033[1;16H\033[42mYOU WIN !!!\033[0m\033[K\033[0m")
	fmt.Printf("\033[24;1H\033[42mYOU WIN !!!\033[0m\033[K\033[0m")
}

func explodeMine(b *board, x int, y int) {
	fmt.Printf("\033[26;1H")
	fmt.Printf("\033[0H\033[2J")
	fmt.Printf("\033[0H\033[1mCurrent Board\n\033[0m")
	printLine(80)
	fmt.Println()
	drawWholeBoard(b, false)
	fmt.Println()
	fmt.Printf("\033[1;16H\033[41mBOOM!!!\033[0m\033[K\033[0m")
	fmt.Printf("\033[24;1H\033[31mYOU LOST !!!\033[0m\033[K\033[0m")
}

func testPositioning(b *board) {
	for y := 1; y <= b.height; y++ {
		for x := 1; x <= b.width; x++ {
			row := ORIGINY + y * 2
			column := ORIGINX + (x-1) * 4 + 2
			fmt.Printf("\033[%d;%dH%s", row, column, showPiece(b, x-1,y-1,false))
//			fmt.Printf("\033[0;50Hrow=%d,column=%d", y, x)
			time.Sleep(50 * time.Millisecond)
		}
	}

	
}
