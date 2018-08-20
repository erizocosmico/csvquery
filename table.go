package csvquery

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"

	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

// Table is a SQL table that will read CSV rows as SQL rows.
type Table struct {
	name   string
	file   string
	schema sql.Schema
	mu     *sync.RWMutex
}

var _ sql.Inserter = (*Table)(nil)

// NewTable creates a new table given a name and filename.
func NewTable(name, file string) (*Table, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	csvr := csv.NewReader(f)
	header, err := csvr.Read()
	if err != nil {
		return nil, fmt.Errorf("csvquery: unable to read header of table %s: %s", name, err)
	}

	t := &Table{name: name, file: file, mu: new(sync.RWMutex)}
	t.schema = make(sql.Schema, len(header))
	for i, col := range header {
		t.schema[i] = &sql.Column{
			Type:     sql.Text,
			Nullable: false,
			Source:   name,
			Name:     col,
			Default:  "",
		}
	}

	return t, nil
}

// Children implements the sql.Table interface.
func (Table) Children() []sql.Node { return nil }

// Name returns the table name.
func (t Table) Name() string { return t.name }

// Resolved implements the sql.Table interface.
func (t Table) Resolved() bool { return true }

// RowIter returns an iterator over all table rows.
func (t Table) RowIter(ctx *sql.Context) (sql.RowIter, error) {
	t.mu.RLock()
	f, err := os.Open(t.file)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		return nil, fmt.Errorf("csvquery: error reading header of table %q: %s", t.name, err)
	}

	return &csvRowIter{closer: f, unlock: t.mu.RUnlock, r: r}, nil
}

// Insert implements the sql.Inserter interface.
func (t Table) Insert(ctx *sql.Context, row sql.Row) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	f, err := os.OpenFile(t.file, os.O_APPEND|os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// We need to go to the last character written and check if it's a
	// new line character. If it's not, we add a new line character
	// before writing the row.
	if _, err := f.Seek(-1, io.SeekEnd); err != nil {
		return err
	}

	var last [1]byte
	if _, err := io.ReadFull(f, last[:]); err != nil {
		return err
	}

	if rune(last[0]) != '\n' {
		if _, err := f.Write([]byte{byte('\n')}); err != nil {
			return err
		}
	}

	// Then seek to the end of the file and write the row.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	var record = make([]string, len(row))
	for i, v := range row {
		record[i] = fmt.Sprint(v)
	}

	return w.Write(record)
}

// Schema returns the table schema.
func (t Table) Schema() sql.Schema { return t.schema }

func (t Table) String() string {
	var columns = make([]string, len(t.schema))
	for i, col := range t.schema {
		columns[i] = col.Name
	}

	tp := sql.NewTreePrinter()
	_ = tp.WriteNode("CSVTable(%s)", t.name)
	_ = tp.WriteChildren(columns...)
	return tp.String()
}

// TransformExpressionsUp implements the sql.Table interface.
func (t Table) TransformExpressionsUp(sql.TransformExprFunc) (sql.Node, error) {
	return t, nil
}

// TransformUp implements the sql.Table interface.
func (t Table) TransformUp(fn sql.TransformNodeFunc) (sql.Node, error) {
	return fn(t)
}
