package csvquery

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	db := NewDatabase("foo")

	require.Equal(t, "foo", db.Name())
	require.NoError(t, db.AddTable("superheroes", filepath.Join("_testdata", "superheroes.csv")))
	require.Error(t, db.AddTable("superheroes", filepath.Join("_testdata", "ratings.csv")))
	require.Error(t, db.AddTable("foo", filepath.Join("_testdata", "does_not_exist.csv")))
	require.Len(t, db.Tables(), 1)
}
