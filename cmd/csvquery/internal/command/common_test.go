package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNameFromPath(t *testing.T) {
	testCases := []struct {
		input, name string
	}{
		{"foo/bar.csv", "bar"},
		{"foo/bar.x.csv", "barx"},
		{"foo/bar1_2.csv", "bar1_2"},
	}

	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tt.name, nameFromPath(tt.input))
		})
	}
}

func TestSplitFile(t *testing.T) {
	testCases := []struct {
		input, name, path string
	}{
		{"foo/bar.csv", "bar", "foo/bar.csv"},
		{"foo/bar.csv:foo", "foo", "foo/bar.csv"},
		{"foo/bar.x.csv", "barx", "foo/bar.x.csv"},
		{"foo/bar1_2.csv", "bar1_2", "foo/bar1_2.csv"},
	}

	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			require := require.New(t)

			name, path := splitFile(tt.input)
			require.Equal(tt.name, name)
			require.Equal(tt.path, path)
		})
	}
}
