package csvquery

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

const testCsv = `1,Jake Peralta
2,Amy Santiago
3,Raymond Holt
4,Regina Linetti
`

func TestCsvIter(t *testing.T) {
	require := require.New(t)

	closer := new(fakeCloser)
	r := csv.NewReader(strings.NewReader(testCsv))
	var unlocked bool
	iter := &csvRowIter{unlock: fakeUnlock(&unlocked), closer: closer, r: r}

	rows, err := sql.RowIterToRows(iter)
	require.NoError(err)

	expected := []sql.Row{
		{"1", "Jake Peralta"},
		{"2", "Amy Santiago"},
		{"3", "Raymond Holt"},
		{"4", "Regina Linetti"},
	}

	require.True(closer.closed)
	require.True(unlocked)
	require.Equal(expected, rows)
}

type fakeCloser struct {
	closed bool
}

func (f *fakeCloser) Close() error {
	f.closed = true
	return nil
}

func fakeUnlock(unlocked *bool) func() {
	return func() {
		*unlocked = true
	}
}
