package movezeroes

import (
	"reflect"
	"testing"
)

func TestMoveZeroes(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{
			name: "mixed",
			nums: []int{0, 1, 0, 3, 12},
			want: []int{1, 3, 12, 0, 0},
		},
		{
			name: "all_zero",
			nums: []int{0, 0, 0},
			want: []int{0, 0, 0},
		},
		{
			name: "no_zero",
			nums: []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "leading_zeros",
			nums: []int{0, 0, 1, 2},
			want: []int{1, 2, 0, 0},
		},
		{
			name: "trailing_zero",
			nums: []int{4, 5, 0},
			want: []int{4, 5, 0},
		},
		{
			name: "single_element_zero",
			nums: []int{0},
			want: []int{0},
		},
		{
			name: "single_element_non_zero",
			nums: []int{7},
			want: []int{7},
		},
		{
			name: "alternating",
			nums: []int{0, 1, 0, 2, 0, 3, 0},
			want: []int{1, 2, 3, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := append([]int(nil), tt.nums...)
			MoveZeroes(in)
			if !reflect.DeepEqual(in, tt.want) {
				t.Fatalf("MoveZeroes(%v) => %v, want %v", tt.nums, in, tt.want)
			}

			// Verify non-zero relative order preserved.
			gotNonZero := filterNonZero(in)
			wantNonZero := filterNonZero(tt.nums)
			if !reflect.DeepEqual(gotNonZero, wantNonZero) {
				t.Fatalf("non-zero order changed: got %v, want %v", gotNonZero, wantNonZero)
			}
		})
	}
}

func filterNonZero(nums []int) []int {
	out := make([]int, 0, len(nums))
	for _, v := range nums {
		if v != 0 {
			out = append(out, v)
		}
	}
	return out
}
