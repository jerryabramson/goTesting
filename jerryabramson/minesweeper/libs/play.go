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
		discoveredBoard: make([][]int64, w, h),
		actualBoard:     make([][]int64, w, h),
		mineCount:       m,
		boardSize:       w * h,
		width:           w,
		height:          h,
		seed:            rand.NewSource(time.Now().UnixNano())}
	b.myRandom = rand.New(b.seed)
	for x := 0; x < w; x++ {
		b.discoveredBoard[x] = make([]int64, h)
		b.actualBoard[x] = make([]int64, h)
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			b.discoveredBoard[x][y] = UNKNOWN
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
		if b.actualBoard[spotW][spotH] != MINE {
			b.actualBoard[spotW][spotH] = MINE
			minesSoFar--
			if minesSoFar == 0 {
				fmt.Printf("Done\n")
				setBoardCounts(b)
				return
			}
		}
	}

}

// determine if a space is a mine
func IsMined(b *board, x int, y int) bool {
	return b.actualBoard[x][y] == MINE
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
	checkXoffset := int(-1)
	checkYoffset := int(-1)
	for true {
		if checkXoffset == 0 && checkYoffset == 0 {
			checkXoffset++
			continue
		} else {
			px := x + checkXoffset
			py := y + checkYoffset
			if px >= 0 && px < b.width {
				if py >= 0 && py < b.height {
					pieceAround := b.actualBoard[px][py]
					if pieceAround == MINE {
						cnt++
					}
				}
			}
			checkXoffset++
			if checkXoffset == 2 {
				checkXoffset = -1
				checkYoffset++
				if checkYoffset == 2 {
					break
				}
			}
		}
	}
	return cnt
}

// Reveal either the board discovered so far,
// or the underlying board itself (after win/lose)
func revealBoard(b *board, h bool) {
	fmt.Println()
	for x := 0; x < b.width; x++ {
		fmt.Printf("\t")
		for y := 0; y < b.height; y++ {
			fmt.Printf("%s", showPiece(b, x, y, h))
		}
		fmt.Println()
	}
	fmt.Println()
}

// Go through each point, and set the count of
// mines around it. We skipp points that are mines
// themselves.
func setBoardCounts(b *board) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			if b.actualBoard[x][y] != MINE {
				c := CountsAroundPoint(b, x, y)
				b.actualBoard[x][y] = c
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
	p := myPieces[x][y]
	if p == MINE {
		return "\033[41;1m \033[0m\033[31;1mX\033[0m\033[41;1m \033[0m"
	} else if p == UNKNOWN {
		return "\033[47m ? \033[0m"
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

// main game loop
func Play(b *board) string {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Printf("\033[2J\033[0H\033[1mCurrent Board\n")
	revealBoard(b, true)
	for true {
		fmt.Printf("Please choose a spot to check for mine [origin at 1,1] (y,x): ")
		err := sc.Scan()
		if !err {
			msg := sc.Err()
			if msg == nil {
				msg = io.EOF
			}
			var errMsg strings.Builder
			fmt.Fprintf(&errMsg, "\nI/O Error: %v\n", msg)
			return errMsg.String()
		}

		input := sc.Text()

		dim := strings.Split(input, ",")
		if len(dim) != 2 {
			fmt.Println("Invalid Syntax")
			continue
		}
		var ierr error
		x, ierr := SafeAtoI(dim[0])
		x--
		if ierr != nil {
			fmt.Println(ierr)
		}
		y, ierr := SafeAtoI(dim[1])
		y--
		if ierr != nil {
			fmt.Println(ierr)
		}
		if !validDimension(b, x, y) {
			fmt.Println("** Out of Range")
			continue
		}
		if b.discoveredBoard[x][y] != UNKNOWN {
			fmt.Printf("You have already revealed space %d,%d: value = %s\n",
				x, y, showPiece(b, x, y, true))
			continue
		}
		if IsMined(b, x, y) {
			fmt.Printf("\033[2J\033[0H")
			fmt.Printf("\033[31mBOOM!!!\033[0m\n")
			revealBoard(b, false)
			return "BOOM"
		}
		b.discoveredBoard[x][y] = b.actualBoard[x][y]
		b.spotsTraversedSoFar++
		if b.spotsTraversedSoFar+b.mineCount == b.boardSize {
			fmt.Printf("\033[2J\033[0H")
			fmt.Printf("\033[32mYOU WIN !!!\033[0m\n")
			return "YOU WIN"
		}
		fmt.Printf("\033[2J\033[0H\033[1mCurrent Board\n")
		revealBoard(b, true)
	}
	return "DONE"
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
