// leetcode problem to find median of two integer arrays
package main

import (
    "fmt"
)


func main() {
	n1 := []int{1,2,3}
	n2 := []int{4,5,6}
	retx := findMedianSortedArrays(n1, n2)
	fmt.Printf("Return is '%f'\n", retx)
}


// Algorithm is to build a single array that
// is a sorted copy of the two arrays, concatenated together.
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    returnValue := float64(0) // return defaults to 0
	totalLength := int(len(nums1)+len(nums2)) // total length of concatanated array
	odd := bool((totalLength % 2) != 0) // whether or no the lengh is even or odd (for median)
    midPoint := int(totalLength / 2) // midpoint, rounded down
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
		if (i >= len(nums1)) {
			combinedNums[k] = nums2[j]
			j++
		} else if (j >= len(nums2)) {
			combinedNums[k] = nums1[i]
			i++
		} else if (nums1[i] < nums2[j]) {
			combinedNums[k] = nums1[i]
			i++
		} else {
			combinedNums[k] = nums2[j]
			j++
		}
		k++;
	}
	if odd {
		// If the total length is odd, the median is trivially the
		// 'middle' entry
		returnValue = float64(combinedNums[midPoint])
	} else {
		// Otherwise, divide the midPoint entry with the following
		// entry (indexes start at 0)
        returnValue = (float64(combinedNums[midPoint-1]) +
                       float64(combinedNums[midPoint])) / 2.0
                     
	}
	return returnValue
}
