package minsubsum

import "testing"

func TestMinSubArrayLen(t *testing.T) {
	tests := []struct {
		name   string
		target int
		nums   []int
		want   int
	}{
		{
			name:   "example",
			target: 7,
			nums:   []int{2, 3, 1, 2, 4, 3},
			want:   2, // [4,3]
		},
		{
			name:   "exact_match",
			target: 4,
			nums:   []int{1, 4, 4},
			want:   1,
		},
		{
			name:   "no_solution",
			target: 15,
			nums:   []int{1, 2, 3, 4, 5},
			want:   5, // total sum 15
		},
		{
			name:   "single_element_valid",
			target: 5,
			nums:   []int{5},
			want:   1,
		},
		{
			name:   "single_element_invalid",
			target: 6,
			nums:   []int{5},
			want:   0,
		},
		{
			name:   "many_small_values",
			target: 11,
			nums:   []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 10},
			want:   2, // [1,10]
		},
		{
			name:   "tight_window",
			target: 8,
			nums:   []int{1, 2, 3, 4, 5},
			want:   2, // [3,5]
		},
		{
			name:   "exact_total",
			target: 21,
			nums:   []int{1, 2, 3, 4, 5, 6},
			want:   6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinSubArrayLen(tt.target, tt.nums)
			if got != tt.want {
				t.Fatalf("MinSubArrayLen(%d, %v) = %d, want %d", tt.target, tt.nums, got, tt.want)
			}
		})
	}
}
