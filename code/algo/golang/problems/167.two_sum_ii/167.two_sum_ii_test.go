package twosumii

import "testing"

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{
			name:   "basic",
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{1, 2},
		},
		{
			name:   "with_negatives",
			nums:   []int{-3, -1, 0, 1, 2},
			target: 1,
			want:   []int{2, 5},
		},
		{
			name:   "large_numbers",
			nums:   []int{1, 2, 3, 4, 4},
			target: 8,
			want:   []int{4, 5},
		},
		{
			name:   "minimal_pair",
			nums:   []int{1, 3},
			target: 4,
			want:   []int{1, 2},
		},
		{
			name:   "zeros",
			nums:   []int{0, 0, 3, 4},
			target: 0,
			want:   []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum(tt.nums, tt.target)
			if got == nil {
				t.Fatalf("TwoSum(%v, %d) returned nil, want %v", tt.nums, tt.target, tt.want)
			}
			if len(got) != 2 || got[0] != tt.want[0] || got[1] != tt.want[1] {
				t.Fatalf("TwoSum(%v, %d) = %v, want %v", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
