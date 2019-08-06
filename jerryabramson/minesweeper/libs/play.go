// Simple implementation of MineSwepper:
//     https://en.wikipedia.org/wiki/Minesweeper_(video_game)
package minesweeper


import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"errors"
)

const (
	MINE                        = -1
	UNKNOWN                     = -1000
	POSSIBLE                    = -100
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
	return cnt
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
			}
		}
	}
}


// determine if a provided dimension is valid
func validDimension(b *board, x int, y int) bool {
	return ((x >= 0 && x < b.width) && (y >= 0 && y < b.height))
}

// main game loop
func Play(b *board) string {
	drawWholeBoard(b, true)
	for true {
		x, y, ierr, safeMark := prompt(b)
		clearStatus()
		if (x == -99) {
			break
		}
		if (x == -1 || y == -1) {
			continue
		}
		if (ierr == nil) {
			if !validDimension(b, x, y) {
				outOfRange()
				continue
			} else {
				if (safeMark) {
					if b.discoveredBoard[y][x] == UNKNOWN {
						b.discoveredBoard[y][x] = POSSIBLE
						revealPiece(b, x, y, true)
						statusMsg(b, "Setting piece as possibly mined", COLOR_BLUE_CSI, x, y)
					} else {
						duplicateMove(b, x, y)
					}
					continue
				}
				if b.discoveredBoard[y][x] == POSSIBLE {
					statusMsg(b, "Re-Setting piece to unknown", COLOR_GREEN_CSI, x, y)
					b.discoveredBoard[y][x] = UNKNOWN
					revealPiece(b, x, y, true)
					continue
				}
				if b.discoveredBoard[y][x] == UNKNOWN {
					revealPiece(b, x, y, false)
				} else {
					duplicateMove(b, x, y)
					continue
				}						
				if IsMined(b, x, y) {
					explodeMine(b, x, y)
					return "You lost."
				}
				b.discoveredBoard[y][x] = b.actualBoard[y][x]
				b.spotsTraversedSoFar++
				if b.spotsTraversedSoFar+b.mineCount == b.boardSize {
					win(b, x, y)
					return "YOU WIN"
				}
			}
		}
	}
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
