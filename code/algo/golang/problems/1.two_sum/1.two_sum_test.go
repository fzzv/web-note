package twosum

import "testing"

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		expect []int
	}{
		{
			name:   "basic",
			nums:   []int{2, 7, 11, 15},
			target: 9,
			expect: []int{0, 1},
		},
		{
			name:   "with_duplicates",
			nums:   []int{3, 2, 4},
			target: 6,
			expect: []int{1, 2},
		},
		{
			name:   "same_value_twice",
			nums:   []int{3, 3},
			target: 6,
			expect: []int{0, 1},
		},
		{
			name:   "negatives",
			nums:   []int{-1, -2, -3, -4, -5},
			target: -8,
			expect: []int{2, 4},
		},
		{
			name:   "later_duplicate",
			nums:   []int{1, 5, 1, 5},
			target: 10,
			expect: []int{1, 3},
		},
		{
			name:   "no_solution",
			nums:   []int{1, 2, 3},
			target: 7,
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum(tt.nums, tt.target)

			if tt.expect == nil {
				if got != nil {
					t.Fatalf("expected no solution, got %v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected indices %v, got nil", tt.expect)
			}

			if len(got) != 2 {
				t.Fatalf("expected 2 indices, got %v", got)
			}

			for _, idx := range got {
				if idx < 0 || idx >= len(tt.nums) {
					t.Fatalf("index %d out of range for nums length %d", idx, len(tt.nums))
				}
			}

			if got[0] == got[1] {
				t.Fatalf("indices must refer to two distinct elements: %v", got)
			}

			sum := tt.nums[got[0]] + tt.nums[got[1]]
			if sum != tt.target {
				t.Fatalf("values at indices %v sum to %d, want %d", got, sum, tt.target)
			}

			if !pairMatches(got, tt.expect) {
				t.Fatalf("expected indices %v, got %v", tt.expect, got)
			}
		})
	}
}

func pairMatches(got, expect []int) bool {
	return (got[0] == expect[0] && got[1] == expect[1]) ||
		(got[0] == expect[1] && got[1] == expect[0])
}
