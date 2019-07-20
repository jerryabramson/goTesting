package main
import (
	"fmt"
	"strings"
	"strconv"
	"math/rand"
	"os"
	"time"
	"bufio"
	"io"
)

const (
    MINE = -1
    UNKNOWN = -1000
)


type board struct {
	discoveredBoard [][]int64
	actualBoard [][]int64
	mineCount int
	boardSize int
	width int
	height int
	seed rand.Source
	myRandom *rand.Rand
	spotsTraversedSoFar int
}

func New(m int, w int, h int) board {
	b := board{
		discoveredBoard : make([][]int64, w, h),
		actualBoard : make ([][]int64, w, h),
        mineCount : m,
        boardSize : w * h,
		width : w,
		height : h,
		seed : rand.NewSource(time.Now().UnixNano())}
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

func PopulateBoard(b *board) {
	minesSoFar := b.mineCount
	for {
		spotW := b.myRandom.Intn(b.width)
		spotH := b.myRandom.Intn(b.height)
		if (b.actualBoard[spotW][spotH] != MINE) {
			b.actualBoard[spotW][spotH] = MINE
			minesSoFar--
			if (minesSoFar == 0) {
				fmt.Printf("Done\n")
				setBoardCounts(b);
				return
			}
		}
	}

}
    
func IsMined(b *board, x int, y int)  bool {
	return b.actualBoard[x][y] == MINE
}
    
func countsAroundPoint(b *board, x int, y int) int64 {
        cnt := int64(0)
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
		checkXoffset := int(-1)
		checkYoffset := int(-1)
		for true {
			if (checkXoffset == 0 && checkYoffset == 0) {
				checkXoffset++;
				continue;
			} else {
				px := x + checkXoffset
				py := y + checkYoffset
				if (px >= 0 && px < b.width) {
					if (py >= 0 && py < b.height) {
						pieceAround := b.actualBoard[px][py]
						if (pieceAround == MINE) {
							cnt++
						}
					}
				}
				checkXoffset++;
				if (checkXoffset == 2) {
					checkXoffset = -1
					checkYoffset++
					if (checkYoffset == 2) {
						break
					}
				}
			}
		}
		return cnt
    }

    func revealBoard(b *board, h bool) {
		// if (h) {
		// 	fmt.Println("Board so far")
		// } else {
		// 	fmt.Println("Actual Board")
		// }
		for x := 0; x < b.width; x++ {
           for y := 0; y < b.height; y++ {
               fmt.Printf("%s", showPiece(b, x,y, h))
           }
           fmt.Println()
		}
		fmt.Println()
    }

 
    func setBoardCounts(b *board) {
        for x := 0; x < b.width; x++ {
            for y := 0; y < b.height; y++ {
                if (b.actualBoard[x][y] != MINE) {
                    c := countsAroundPoint(b, x,y)
//					fmt.Printf("point %d,%d : count = %d\n", x+1,y+1,c)
                    b.actualBoard[x][y] = c
                } else {
//					fmt.Printf("Mine at %d,%d\n", x+1,y+1)
				}
            }
        }
    }

func showPiece(b *board, x int, y int, h bool) string {

        myPieces := b.discoveredBoard
        if (!h) {
            myPieces = b.actualBoard
        }
	    p := myPieces[x][y]
        if (p == MINE) {
            return "\033[41;1m \033[0m\033[31;1m*\033[0m\033[41;1m \033[0m"
        } else if (p == UNKNOWN) {
            return "\033[47m ? \033[0m"
        } else {
            ps := strconv.FormatInt(p, 10)
			pString := string("X")
            if (p > 0) {
                pString = "\033[42m " + ps + " \033[0m"
            } else {
				pString = "\033[46m " + ps + " \033[0m"
			}
            return pString
        }
    }

    func validDimension(b *board, x int, y int) bool {
        return ((x >= 0 && x < b.width) && (y >= 0 && y < b.height))
    }

func Play(b *board) {
		sc := bufio.NewScanner(os.Stdin)
            fmt.Printf("\033[2J\033[0H\033[1mCurrent Board\n")
        revealBoard(b, true)
        for true {
			fmt.Printf("Please choose a spot to check for mine [origin at 1,1] (y,x): ")
			err := sc.Scan()
			if (!err) {
				msg := sc.Err()
				if (msg == nil) {
					msg = io.EOF
				}
				fmt.Printf("\nI/O Error: %v\n", msg)
				os.Exit(0)
			}

            input := sc.Text()

            dim := strings.Split(input, ",")
            if (len(dim) != 2) {
                fmt.Println("Invalid Syntax")
                continue
            } 
            x := SafeAtoI(dim[0]) - 1
            y := SafeAtoI(dim[1]) - 1
            if (!validDimension(b, x,y)) {
                fmt.Println("** Out of Range");
                continue
            }
            if (b.discoveredBoard[x][y] != UNKNOWN) {
                fmt.Printf("You have already revealed space %d,%d: value = %s\n", 
					x,y, showPiece(b, x,y, true))
				continue
            }
            if (IsMined(b, x,y)) {
            fmt.Printf("\033[2J\033[0H")
                fmt.Printf("\033[31mBOOM!!!\033[0m\n")
                revealBoard(b, false)
                return
            }
            b.discoveredBoard[x][y] = b.actualBoard[x][y]
			b.spotsTraversedSoFar++;
			if (b.spotsTraversedSoFar + b.mineCount == b.boardSize) {
            fmt.Printf("\033[2J\033[0H")
				fmt.Printf("\033[32mYOU WIN !!!\033[0m\n")
                revealBoard(b, false)
				os.Exit(0)
			}
            fmt.Printf("\033[2J\033[0H\033[1mCurrent Board\n")
            revealBoard(b, true)
        }
    }

