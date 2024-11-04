package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckLeadingZeroBits(t *testing.T) {
	tests := []struct {
		leadingZeroBits uint
		hash            [32]byte
		expected        bool
	}{
		{
			leadingZeroBits: 0,
			hash:            [32]byte{255}, // 11111111 00000000 ...
			expected:        true,
		},
		{
			leadingZeroBits: 1,
			hash:            [32]byte{255}, // 11111111 00000000 ...
			expected:        false,
		},
		{
			leadingZeroBits: 1,
			hash:            [32]byte{128}, // 10000000 00000000 ...
			expected:        false,
		},
		{
			leadingZeroBits: 1,
			hash:            [32]byte{64}, // 01000000 00000000 ...
			expected:        true,
		},
		{
			leadingZeroBits: 2,
			hash:            [32]byte{64}, // 01000000 00000000 ...
			expected:        false,
		},
		{
			leadingZeroBits: 2,
			hash:            [32]byte{63}, // 00111111 00000000 ...
			expected:        true,
		},
		{
			leadingZeroBits: 1000,        // just check that all 256 bits are zeroes
			hash:            [32]byte{0}, // 00000000 00000000 ...
			expected:        true,
		},
	}

	for i, test := range tests {
		res := CheckLeadingZeroBits(test.leadingZeroBits, test.hash)
		require.Equal(t, test.expected, res, fmt.Sprintf("case %d failed", i))
	}
}
