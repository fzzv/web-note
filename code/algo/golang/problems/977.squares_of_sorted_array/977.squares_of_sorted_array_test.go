package sortedsquares

import (
	"reflect"
	"testing"
)

func TestSortedSquares(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{
			name: "mixed_neg_pos",
			nums: []int{-4, -1, 0, 3, 10},
			want: []int{0, 1, 9, 16, 100},
		},
		{
			name: "all_negative",
			nums: []int{-7, -3, -1},
			want: []int{1, 9, 49},
		},
		{
			name: "all_positive",
			nums: []int{1, 2, 3},
			want: []int{1, 4, 9},
		},
		{
			name: "single_zero",
			nums: []int{0},
			want: []int{0},
		},
		{
			name: "two_elements",
			nums: []int{-2, 3},
			want: []int{4, 9},
		},
		{
			name: "duplicates",
			nums: []int{-2, -2, 2, 2},
			want: []int{4, 4, 4, 4},
		},
		{
			name: "edge_order",
			nums: []int{-5, -3, -2, -1},
			want: []int{1, 4, 9, 25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortedSquares(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SortedSquares(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
