// Simple implementation of MineSwepper:
//     https://en.wikipedia.org/wiki/Minesweeper_(video_game)
package minesweeper


import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"errors"
)

// constants - MINE is a mine
const (
	MINE    = -1
	UNKNOWN = -1000
	POSSIBLE = 100
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
)

// encapsulate the board
type board struct {
	discoveredBoard     [][]int64
	actualBoard         [][]int64
	mineCount           int
	boardSize           int
	width               int
	height              int
	seed                rand.Source
	myRandom            *rand.Rand
	spotsTraversedSoFar int
}

// Create a new board. Note, Go has no real two-dimensional
// arrays. Instead, you have an array of arrays.
func New(m int, w int, h int) board {
	b := board{
		discoveredBoard: make([][]int64, w*h),
		actualBoard:     make([][]int64, w*h),
		mineCount:       m,
		boardSize:       w * h,
		width:           w,
		height:          h,
		seed:            rand.NewSource(time.Now().UnixNano())}
	b.myRandom = rand.New(b.seed)
	for y := 0; y < h; y++ {
		b.discoveredBoard[y] = make([]int64, w)
		b.actualBoard[y] = make([]int64, w)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			b.discoveredBoard[y][x] = UNKNOWN
		}
	}

	return b
}

// place specified number of mines throughout the board
func PopulateBoard(b *board) {
	minesSoFar := b.mineCount
	for {
		spotW := b.myRandom.Intn(b.width)
		spotH := b.myRandom.Intn(b.height)
		if b.actualBoard[spotH][spotW] != MINE {
			b.actualBoard[spotH][spotW] = MINE
			minesSoFar--
			if minesSoFar == 0 {
//				fmt.Printf("Done\n")
				setBoardCounts(b)
				return
			}
		}
	}

}

// determine if a space is a mine
func IsMined(b *board, x int, y int) bool {
	return b.actualBoard[y][x] == MINE
}

// Count the number of mines surrounding a given point
// on the board
	/*
	 * For a given piece, 'p', we count the number
	 * of mines, as show below.
	 *
	 * So, we start with x-1, and go to x+1.
	 * And y-1, to y+1
	 * If x or y is less than zero, skip.
	 * if x or y is equal to length, skip.
	 *       0   1   2
	 *    +-------------+
	 *    |             |
	 *  0 |  x   x   x  |
	 *    |             |
	 *  1 |  x   p   x  |
	 *    |             |
	 *  2 |  x   x   x  |
	 *    |             |
	 *    +-------------+
	 */
func CountsAroundPoint(b *board, x int, y int) int64 {
	cnt := int64(0)
//	fmt.Printf("piecesAround y=%d,x=%d:\n", y+1,x+1)
	checkXoffset := []int {-1,0,1}
	checkYoffset := []int {-1,0,1}
	var px int
	var py int
	for iy:= 0; iy < len(checkYoffset); iy++ {
		py = y + checkYoffset[iy]
		for ix := 0; ix < len(checkXoffset); ix++ {
			px = x + checkXoffset[ix]
			if (checkXoffset[ix] == 0 && checkYoffset[iy] == 0) {
//				fmt.Printf("Origin")
			} else if px >= 0 && px < b.width && py >= 0 && py < b.height {
//				fmt.Printf("  check py=%d,px=%d; ", py+1,px+1)
				pieceAround := b.actualBoard[py][px]
//				fmt.Printf("piece = %d", pieceAround)
				if pieceAround == MINE {
					cnt++
				}
			} else {
//				fmt.Printf("Out of bounds at py=%d,px=%d", py+1,px+1)
			}
//			fmt.Println()
		}
	}
//	fmt.Printf("count = %d\n", cnt)
	return cnt
}

// Reveal either the board discovered so far,
// or the underlying board itself (after win/lose)
func revealBoard(b *board, h bool) {
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

// Go through each point, and set the count of
// mines around it. We skipp points that are mines
// themselves.
func setBoardCounts(b *board) {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if b.actualBoard[y][x] != MINE {
				c := CountsAroundPoint(b, x, y)
				b.actualBoard[y][x] = c
			} else {
//				fmt.Printf("mine at y=%d,x=%d\n", y+1,x+1)
			}
		}
	}
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

// determine if a provided dimension is valid
func validDimension(b *board, x int, y int) bool {
	return ((x >= 0 && x < b.width) && (y >= 0 && y < b.height))
}

func revealPiece(b *board, x int, y int, h bool) {
	originY := 4
	originX := 8
	row := originY + y * 2 + 1
	column := originX + (x-1) * 4 + 2
//	fmt.Printf("\033[0;50Hrow=%d,column=%d", y, x)
	fmt.Printf("\033[%d;%dH%s", row, column, showPiece(b, x,y, h))
}

func testPositioning(b *board) {
	originY := 4
	originX := 8
	for y := 1; y <= b.height; y++ {
		for x := 1; x <= b.width; x++ {
			row := originY + y * 2
			column := originX + (x-1) * 4 + 2
			fmt.Printf("\033[%d;%dH%s", row, column, showPiece(b, x-1,y-1,false))
//			fmt.Printf("\033[0;50Hrow=%d,column=%d", y, x)
			time.Sleep(50 * time.Millisecond)
		}
	}

	
}
// main game loop
func Play(b *board) string {
	sc := bufio.NewScanner(os.Stdin)

	fmt.Printf("\033[0H\033[2J")
	fmt.Printf("\033[0H\033[1mCurrent Board\n\033[0m")
	printLine(80)
	fmt.Println()
	revealBoard(b, true)
	fmt.Println()
//	testPositioning(b)
	
	
	for true {
		fmt.Printf("\033[25;1HPlease choose a spot to check for a mine [origin at 1,1] (y,x)[,?]:\033[K ")
		err := sc.Scan()
		if !err {
			msg := sc.Err()
			if msg == nil {
				msg = io.EOF
			}
			var errMsg strings.Builder
			fmt.Fprintf(&errMsg, "\033[1;16H\033[31mI/O Error: %v\033[K\033[m", msg)
			fmt.Printf("\033[26;1H")
			return errMsg.String()
		}

		input := sc.Text()
		dim := strings.Split(input, ",")
		safeMark := bool(false)
		if len(dim) == 3 {
			if strings.Compare(dim[2], "?") == 0 {
				safeMark = true;
			}
		} else if len(dim) != 2 {
			fmt.Printf("\033[1;16H\033[31mInvalid Syntax\033[K\033[0m")
			continue
		}
		var ierr error
		y, ierr := SafeAtoI(dim[0])
		y--
		if ierr != nil {
			fmt.Print("\033[1;16H\033[31m%s\033[K\033[0m", ierr)
			continue
		}
		x, ierr := SafeAtoI(dim[1])
		x--
		if ierr != nil {
			fmt.Printf("\033[1;16H\033[31m%s\033[K\033[0m", ierr)
			continue
		}

		if !validDimension(b, x, y) {
			fmt.Printf("\033[1;16H\033[31m** Out of Range\033[K\033[0m")
			continue
		}

		if (!safeMark) {
			revealPiece(b, x, y, false)
		}
		if b.discoveredBoard[y][x] != UNKNOWN {
			if b.discoveredBoard[y][x] == POSSIBLE {
				fmt.Printf("\033[1;16H\033[32mresetting space %d,%d from\033[0m %s ",
					x+1, y+1, showPiece(b, x, y, true))
				b.discoveredBoard[y][x] = UNKNOWN
				fmt.Printf(" \033[32mback to\033[0m %s\033[K", showPiece(b,x,y,true))
				revealPiece(b, x, y, true)
			} else {
				fmt.Printf("\033[1;16H\033[33mYou have already revealed space %d,%d: value = \033[0m%s\033[K",
					x+1, y+1, showPiece(b, x, y, true))
			}
			continue
		}
		
		if (safeMark) {
			if (b.discoveredBoard[y][x] != UNKNOWN) {
				fmt.Printf("\033[1;16H\033[33mYou have already revealed space %d,%d: value = \033[0m%s\033[K",
					y+1, x+1, showPiece(b, x, y, true))
			} else {				
				b.discoveredBoard[y][x] = POSSIBLE
				revealPiece(b, x, y, true)
			}
			continue
		} else if IsMined(b, x, y) {
//			revealBoard(b, false)
			fmt.Printf("\033[26;1H")
			fmt.Printf("\033[0H\033[2J")
			fmt.Printf("\033[0H\033[1mCurrent Board\n\033[0m")
			printLine(80)
			fmt.Println()
			revealBoard(b, false)
			fmt.Println()
			fmt.Printf("\033[1;16H\033[41mBOOM!!!\033[0m\033[K\033[0m")
			fmt.Printf("\033[26;1H\033[31mYOU LOST !!!\033[0m\033[K\033[0m\n")
			return "You lost."
		} else {
			b.discoveredBoard[y][x] = b.actualBoard[y][x]
			b.spotsTraversedSoFar++
			if b.spotsTraversedSoFar+b.mineCount == b.boardSize {
//				fmt.Printf("\033[2J\033[0H")
			fmt.Printf("\033[26;1H")
			fmt.Printf("\033[0H\033[2J")
			fmt.Printf("\033[0H\033[1mCurrent Board\n\033[0m")
			printLine(80)
			fmt.Println()
			revealBoard(b, false)
			fmt.Println()
				fmt.Printf("\033[1;16H\033[42mYOU WIN !!!\033[0m\033[K\033[0m")
				fmt.Printf("\033[26;1H\033[42mYOU WIN !!!\033[0m\033[K\033[0m\n")
				return "YOU WIN"
			}
		}
		fmt.Printf("\033[1;16H\033[K")
//		printLine(80)
//		fmt.Println()
//		fmt.Printf("\033[1mCurrent Board\n")
//		revealBoard(b, true)

	}
	fmt.Printf("\033[26;1H")
	return "End of loop."
}

func printLine(n int) {
	for x := 0 ; x < n; x++ {
		fmt.Printf(ANSI_LINE)
	}
}



// utility to safely convert a string to an int, with proper
// error reporting.


// global err
var (
	err error
)

// Input is a string, returns a value, and an error code
func SafeAtoI(val string) (int, error) {
	ret, e := strconv.Atoi(val)
	if e == nil {
       // no errors 
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
