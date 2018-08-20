package csvquery

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

var expectedRatings = []sql.Row{
	{"ww", "foo_bar", "7"},
	{"bm", "foo_bar", "6"},
	{"sm", "foo_bar", "5"},
	{"ma", "john_doe", "8"},
	{"bm", "john_doe", "7"},
	{"ww", "alice_doe", "8"},
	{"ma", "alice_doe", "9"},
}

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

	require.Equal(expectedRatings, rows)
}

func TestInsert(t *testing.T) {
	require := require.New(t)
	file, cleanup := testCSV(t)
	defer cleanup()

	table, err := NewTable("ratings", file)
	require.NoError(err)

	row := sql.NewRow("ma", "miss_america_fan", "10")
	require.NoError(table.Insert(sql.NewEmptyContext(), row))

	iter, err := table.RowIter(sql.NewEmptyContext())
	require.NoError(err)

	rows, err := sql.RowIterToRows(iter)
	require.NoError(err)

	expected := append(expectedRatings, row)
	require.Equal(expected, rows)
}

func TestInsertConcurrent(t *testing.T) {
	require := require.New(t)
	file, cleanup := testCSV(t)
	defer cleanup()

	table, err := NewTable("ratings", file)
	require.NoError(err)

	var iters []sql.RowIter
	for i := 0; i < 5; i++ {
		iter, err := table.RowIter(sql.NewEmptyContext())
		require.NoError(err)

		iters = append(iters, iter)
	}

	var reads int
	go func() {
		time.Sleep(50 * time.Millisecond)
		for _, iter := range iters {
			rows, err := sql.RowIterToRows(iter)
			require.NoError(err)
			require.Equal(expectedRatings, rows)
			reads++
		}
	}()

	row := sql.NewRow("ma", "miss_america_fan", "10")
	require.NoError(table.Insert(sql.NewEmptyContext(), row))

	require.Equal(reads, len(iters))

	iter, err := table.RowIter(sql.NewEmptyContext())
	require.NoError(err)

	rows, err := sql.RowIterToRows(iter)
	require.NoError(err)

	expected := append(expectedRatings, row)
	require.Equal(expected, rows)
}

func testCSV(t *testing.T) (file string, cleanup func()) {
	t.Helper()
	content, err := ioutil.ReadFile(filepath.Join("_testdata", "ratings.csv"))
	require.NoError(t, err)

	f, err := ioutil.TempFile(os.TempDir(), "csvquery_")
	require.NoError(t, err)

	_, err = f.Write(content)
	require.NoError(t, err)

	require.NoError(t, f.Close())

	return f.Name(), func() {
		require.NoError(t, os.Remove(f.Name()))
	}
}
