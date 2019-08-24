package edac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryInfoHasErrors(t *testing.T) {
	tests := []struct {
		name     string
		mi       MemoryInfo
		expected bool
	}{
		{
			name: "no-errors",
		},
		{
			name:     "correctable",
			mi:       MemoryInfo{Correctable: 1},
			expected: true,
		},
		{
			name:     "correctable-noinfo",
			mi:       MemoryInfo{CorrectableNoInfo: 1},
			expected: true,
		},
		{
			name:     "uncorrectable",
			mi:       MemoryInfo{Uncorrectable: 1},
			expected: true,
		},
		{
			name:     "uncorrectable-noinfo",
			mi:       MemoryInfo{UncorrectableNoInfo: 1},
			expected: true,
		},
		{
			name: "all",
			mi: MemoryInfo{
				Correctable:         1,
				CorrectableNoInfo:   1,
				Uncorrectable:       1,
				UncorrectableNoInfo: 1,
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mi.Name = "mc0"
			require.Equal(t, tc.expected, tc.mi.HasErrors())
			require.Contains(t, tc.mi.String(), tc.mi.Name)
		})
	}
}

func TestDimmRankHasErrors(t *testing.T) {
	tests := []struct {
		name     string
		dr       DimmRank
		expected bool
	}{
		{
			name: "no-errors",
		},
		{
			name:     "correctable",
			dr:       DimmRank{Correctable: 1},
			expected: true,
		},
		{
			name:     "uncorrectable",
			dr:       DimmRank{Uncorrectable: 1},
			expected: true,
		},
		{
			name: "all",
			dr: DimmRank{
				Correctable:   1,
				Uncorrectable: 1,
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.dr.Name = "dimm0"
			require.Equal(t, tc.expected, tc.dr.HasErrors())
			require.Contains(t, tc.dr.String(), tc.dr.Name)
		})
	}
}
