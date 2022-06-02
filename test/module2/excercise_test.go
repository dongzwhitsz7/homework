package module2

import (
	"dongzw/dongzwhom/module2"
	"fmt"
	"testing"
)

func qs(nums []int, left, right int) []int{
	if left >= right {
		return nums
	}
	p := partition(nums, left, right)
	qs(nums, left, p - 1)
	qs(nums, p + 1, right)
	return nums
}

func partition(nums []int, left, right int) int{
	pivot := nums[left]
	index := left + 1
	for i := index + 1; i <= right; i++ {
		if nums[i] <= pivot {
			nums[i], nums[index] = nums[index], nums[i]
			index++
		}
	}
	x := ""
	for _, v := range x {
		v
	}
	nums[index] = pivot
	return index
}
func TestExcersizeHttpServer(t *testing.T) {
	a := []int{-1,0,1,2,-1,-4}
	b := partition(a, 0, len(a) - 1)
	fmt.Println( a, b)
	server := module2.NewHttpServer()
	server.Serve()
}
