package calc

import "errors"

func Sum(nums []int) int {
	t := 0
	for _, n := range nums {
		t += n
	}
	return t
}

var ErrEmpty = errors.New("empty slice")

func Average(nums []int) (float64, error) {
	if len(nums) == 0 {
		return 0, ErrEmpty
	}
	s := 0
	for _, n := range nums {
		s += n
	}
	return float64(s) / float64(len(nums)), nil
} 
