package reversestring

import (
	"reflect"
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want []byte
	}{
		{
			name: "basic_even",
			in:   []byte{'h', 'e', 'l', 'l', 'o', '!'},
			want: []byte{'!', 'o', 'l', 'l', 'e', 'h'},
		},
		{
			name: "basic_odd",
			in:   []byte{'H', 'a', 'n', 'n', 'a', 'h'},
			want: []byte{'h', 'a', 'n', 'n', 'a', 'H'},
		},
		{
			name: "single_char",
			in:   []byte{'a'},
			want: []byte{'a'},
		},
		{
			name: "repeated_chars",
			in:   []byte{'a', 'a', 'a', 'a'},
			want: []byte{'a', 'a', 'a', 'a'},
		},
		{
			name: "mixed_symbols",
			in:   []byte{'1', '2', '#', ' ', 'b'},
			want: []byte{'b', ' ', '#', '2', '1'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCopy := append([]byte(nil), tt.in...)
			ReverseString(inCopy)
			if !reflect.DeepEqual(inCopy, tt.want) {
				t.Fatalf("ReverseString(%v) = %v, want %v", tt.in, inCopy, tt.want)
			}
		})
	}
}
