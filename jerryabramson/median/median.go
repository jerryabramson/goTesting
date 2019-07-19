// leetcode problem to find median of two integer arrays
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func Log(d bool, f string, msg ...interface{}) {
	if d {
		fmt.Printf(f, msg...)
	}
}

func main() {
	testCases := int(10)
	err := errors.New("none")
	dbg := bool(false)
	for argc := 1; argc < len(os.Args); argc++ {
		a := os.Args[argc]
		if strings.Compare(a, "-d") == 0 {
			dbg = true
		} else {
			testCases, err = strconv.Atoi(os.Args[argc])
			if err != nil {
				fmt.Printf("Error %v\n", err)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Running %d tests\n", testCases)
	t := make([]time.Duration, testCases)
	for caseNumber := 0; caseNumber < testCases; caseNumber++ {
		sz := rand.Intn(1000) + 1 // size of each array to sort
		n1 := make([]int, sz)
		n2 := make([]int, sz)
		for ind := int(0); ind < sz; ind++ {
			n1[ind] = rand.Intn(999) // random number between 1 and 999
			n2[ind] = rand.Intn(999) // random number between 1 and 999
		}
		Log(dbg, "\nTest Case %2d: Arrays of size %3d\n", caseNumber+1, sz)
		if dbg {
			PrintArray(n1)
			PrintArray(n2)
		}
		t0 := time.Now()
		retx := float64(0)
		if dbg {
			retx = FindMedianSortedArrays(n1, n2)
			Log(dbg, "Median = %.2f ", retx)
		} else {
			FindMedianSortedArrays(n1, n2)
		}
		t1 := time.Now()
		t[caseNumber] = t1.Sub(t0)
		Log(dbg, "(%v secs)\n", t1.Sub(t0))
	}
	a, m := FindAverage(t)
	fmt.Printf("\nAverage execution time %v seconds (max %v seconds)\n\n\n",
		a, m)

}

func FindAverage(d []time.Duration) (time.Duration, time.Duration) {
	sum := time.Duration(0)
	max := time.Duration(0)
	for i := 0; i < len(d); i++ {
		sum = sum + d[i]
		if d[i] > max {
			max = d[i]
		}
	}
	avg := time.Duration(int64(sum) / int64(len(d)))
	return avg, max
}

func PrintArray(a []int) {
	fmt.Printf("\tArray: [")
	i := int(0)
	for i < len(a)-1 {
		fmt.Printf("%3d, ", a[i])
		i++
		if i > 10 {
			fmt.Printf("... ")
			break
		}
	}
	fmt.Printf("%d]\n", a[i])
}

// Algorithm is to build a single array that
// is a sorted copy of the two arrays, concatenated together.
func FindMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	returnValue := float64(0)                   // return defaults to 0
	totalLength := int(len(nums1) + len(nums2)) // total length of concatanated array
	odd := bool((totalLength % 2) != 0)         // whether or no the lengh is even or odd (for median)
	midPoint := int(totalLength / 2)            // midpoint, rounded down
	//    fmt.Printf("totalLentgh = '%d', midPoint= '%d'\n", totalLength, midPoint) // debug stmt
	combinedNums := make([]int, totalLength) // make a new array slice of totalLengh size
	// init counters
	i := int(0)
	j := int(0)
	k := int(0)
	// Loop over both int arrays, inserting lowest value into
	// combined array. If one of the int arrays is exhausted, simply
	// insert the next entry from the other int array
	for i < len(nums1) || j < len(nums2) {
		if i >= len(nums1) {
			combinedNums[k] = nums2[j]
			j++
		} else if j >= len(nums2) {
			combinedNums[k] = nums1[i]
			i++
		} else if nums1[i] < nums2[j] {
			combinedNums[k] = nums1[i]
			i++
		} else {
			combinedNums[k] = nums2[j]
			j++
		}
		k++
	}
	if odd {
		// If the total length is odd, the median is trivially the 'middle' entry
		returnValue = float64(combinedNums[midPoint])
	} else {
		// Otherwise, divide the midPoint entry with the following entry (indexes start at 0)
		returnValue = (float64(combinedNums[midPoint-1]) + float64(combinedNums[midPoint])) / 2.0
	}
	return returnValue
}
