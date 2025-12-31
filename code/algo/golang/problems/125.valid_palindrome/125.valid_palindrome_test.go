package validpalindrome

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "classic_true",
			s:    "A man, a plan, a canal: Panama",
			want: true,
		},
		{
			name: "classic_false",
			s:    "race a car",
			want: false,
		},
		{
			name: "digits_and_letters_mismatch",
			s:    "0P",
			want: false,
		},
		{
			name: "only_punctuation",
			s:    ".,,",
			want: true,
		},
		{
			name: "spaces_only",
			s:    "     ",
			want: true,
		},
		{
			name: "mixed_with_symbols_true",
			s:    "ab@a",
			want: true,
		},
		{
			name: "single_char",
			s:    "x",
			want: true,
		},
		{
			name: "empty",
			s:    "",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.s)
			if got != tt.want {
				t.Fatalf("IsPalindrome(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
