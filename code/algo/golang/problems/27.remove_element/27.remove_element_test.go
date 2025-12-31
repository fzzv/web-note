package removeelement

import (
	"reflect"
	"testing"
)

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		val    int
		expect []int
		wantN  int
	}{
		{
			name:   "basic_duplicates",
			nums:   []int{3, 2, 2, 3},
			val:    3,
			expect: []int{2, 2},
			wantN:  2,
		},
		{
			name:   "no_removal",
			nums:   []int{1, 2, 3},
			val:    4,
			expect: []int{1, 2, 3},
			wantN:  3,
		},
		{
			name:   "all_removed",
			nums:   []int{1, 1, 1},
			val:    1,
			expect: []int{},
			wantN:  0,
		},
		{
			name:   "interleaved_values",
			nums:   []int{0, 1, 2, 2, 3, 0, 4, 2},
			val:    2,
			expect: []int{0, 1, 3, 0, 4},
			wantN:  5,
		},
		{
			name:   "single_element_removed",
			nums:   []int{5},
			val:    5,
			expect: []int{},
			wantN:  0,
		},
		{
			name:   "single_element_kept",
			nums:   []int{5},
			val:    1,
			expect: []int{5},
			wantN:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := append([]int(nil), tt.nums...)
			n := RemoveElement(input, tt.val)

			if n != tt.wantN {
				t.Fatalf("RemoveElement(%v, %d) returned length %d, want %d", tt.nums, tt.val, n, tt.wantN)
			}

			if !reflect.DeepEqual(input[:n], tt.expect) {
				t.Fatalf("after RemoveElement, first %d elements = %v, want %v", n, input[:n], tt.expect)
			}

			// Ensure removed value does not appear in the kept prefix.
			for i := 0; i < n; i++ {
				if input[i] == tt.val {
					t.Fatalf("value %d should have been removed but found at index %d", tt.val, i)
				}
			}
		})
	}
}
