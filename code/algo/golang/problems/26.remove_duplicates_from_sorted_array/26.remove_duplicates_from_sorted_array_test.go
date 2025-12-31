package removeduplicates

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		expect []int
		wantN  int
	}{
		{
			name:   "basic_duplicates",
			nums:   []int{1, 1, 2},
			expect: []int{1, 2},
			wantN:  2,
		},
		{
			name:   "all_unique",
			nums:   []int{0, 1, 2, 3},
			expect: []int{0, 1, 2, 3},
			wantN:  4,
		},
		{
			name:   "all_same",
			nums:   []int{5, 5, 5, 5},
			expect: []int{5},
			wantN:  1,
		},
		{
			name:   "empty",
			nums:   []int{},
			expect: []int{},
			wantN:  0,
		},
		{
			name:   "single_element",
			nums:   []int{9},
			expect: []int{9},
			wantN:  1,
		},
		{
			name:   "tail_duplicates",
			nums:   []int{0, 0, 1, 1, 1, 2, 3, 3, 4},
			expect: []int{0, 1, 2, 3, 4},
			wantN:  5,
		},
		{
			name:   "alternating_pairs",
			nums:   []int{1, 1, 2, 2, 3, 3},
			expect: []int{1, 2, 3},
			wantN:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := append([]int(nil), tt.nums...)
			n := RemoveDuplicates(in)

			if n != tt.wantN {
				t.Fatalf("RemoveDuplicates(%v) returned length %d, want %d", tt.nums, n, tt.wantN)
			}

			if !reflect.DeepEqual(in[:n], tt.expect) {
				t.Fatalf("after RemoveDuplicates, prefix = %v, want %v", in[:n], tt.expect)
			}

			// Ensure prefix strictly increasing (no duplicates remain).
			for i := 1; i < n; i++ {
				if in[i] == in[i-1] {
					t.Fatalf("duplicate value %d remains at indices %d and %d", in[i], i-1, i)
				}
			}
		})
	}
}
