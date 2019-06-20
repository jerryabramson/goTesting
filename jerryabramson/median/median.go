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



func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
 
    returnValue := float64(0)
	totalLength := int(len(nums1)+len(nums2))
	odd := bool((totalLength % 2) != 0)
    midPoint := int(totalLength / 2)
    fmt.Printf("totalLentgh = '%d', midPoint= '%d'\n", totalLength, midPoint)
	combinedNums := make([]int, totalLength)
	i := int(0)
	j := int(0)
	k := int(0)
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
		returnValue = float64(combinedNums[midPoint])
	} else {
        returnValue = (float64(combinedNums[midPoint-1]) +
                       float64(combinedNums[midPoint])) / 2.0
                     
	}
	return returnValue
}
