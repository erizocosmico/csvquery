package csvquery

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

func TestTable(t *testing.T) {
	require := require.New(t)

	table, err := NewTable("ratings", filepath.Join("_testdata", "ratings.csv"))
	require.NoError(err)

	require.Equal("ratings", table.Name())
	require.Equal(
		sql.Schema{
			{Name: "superhero_id", Type: sql.Text, Default: "", Source: "ratings"},
			{Name: "username", Type: sql.Text, Default: "", Source: "ratings"},
			{Name: "rating", Type: sql.Text, Default: "", Source: "ratings"},
		},
		table.Schema(),
	)

	iter, err := table.RowIter(sql.NewEmptyContext())
	require.NoError(err)

	rows, err := sql.RowIterToRows(iter)
	require.NoError(err)

	expected := []sql.Row{
		{"ww", "foo_bar", "7"},
		{"bm", "foo_bar", "6"},
		{"sm", "foo_bar", "5"},
		{"ma", "john_doe", "8"},
		{"bm", "john_doe", "7"},
		{"ww", "alice_doe", "8"},
		{"ma", "alice_doe", "9"},
	}
	require.Equal(expected, rows)
}
