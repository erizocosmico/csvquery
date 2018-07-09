package csvquery

import (
	"fmt"

	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

// Database that contains all CSV files as tables. Adding and reading tables
// from the database is not thread-safe. Adding tables should happen before
// they are going to be read.
type Database struct {
	name   string
	tables map[string]sql.Table
}

// NewDatabase creates a new database with the given name.
func NewDatabase(name string) *Database {
	return &Database{
		name:   name,
		tables: make(map[string]sql.Table),
	}
}

// Name returns the name of the database.
func (d *Database) Name() string {
	return d.name
}

// Tables returns a map of the tables indexed by name.
func (d *Database) Tables() map[string]sql.Table {
	return d.tables
}

// AddTable adds a new table with the given name and path.
func (d *Database) AddTable(name, path string) error {
	if _, ok := d.tables[name]; ok {
		return fmt.Errorf("table with name %q already registered", name)
	}

	t, err := NewTable(name, path)
	if err != nil {
		return fmt.Errorf("unable to add table %q: %s", name, err)
	}

	d.tables[name] = t

	return nil
}
