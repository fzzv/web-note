package sortanarray

import (
	"reflect"
	"testing"
)

func TestSortArray(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{
			name: "already_sorted",
			nums: []int{1, 2, 3, 4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "reverse_sorted",
			nums: []int{5, 4, 3, 2, 1},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "with_duplicates",
			nums: []int{2, 2, 1, 1, 3},
			want: []int{1, 1, 2, 2, 3},
		},
		{
			name: "negatives_and_zero",
			nums: []int{0, -1, 3, -5, 2},
			want: []int{-5, -1, 0, 2, 3},
		},
		{
			name: "single_element",
			nums: []int{42},
			want: []int{42},
		},
		{
			name: "empty",
			nums: []int{},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortArray(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SortArray(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
