package command

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/erizocosmico/csvquery"
	sqle "gopkg.in/src-d/go-mysql-server.v0"
)

type baseCmd struct {
	Name  string   `long:"dbname" short:"d" default:"csv" description:"Database name."`
	Files []string `long:"file" short:"f" description:"Add file as a table. You can use the flag in the format '/path/to/file:NAME' to give the file a specific table name. Otherwise, the file name without extension will be the table name with only alphanumeric characters and underscores."`
}

func (b baseCmd) engine() (*sqle.Engine, *csvquery.Database, error) {
	db := csvquery.NewDatabase(b.Name)

	for _, f := range b.Files {
		name, path := splitFile(f)
		if err := db.AddTable(name, path); err != nil {
			return nil, nil, err
		}
	}

	engine := sqle.NewDefault()
	engine.AddDatabase(db)

	if err := engine.Init(); err != nil {
		return nil, nil, fmt.Errorf("csvquery: unable to initialize engine: %s", err)
	}

	return engine, db, nil
}

func splitFile(s string) (name, path string) {
	if idx := strings.LastIndex(s, ":"); idx >= 0 {
		path = s[:idx]
		name = s[idx+1:]
	} else {
		path = s
	}

	if name == "" {
		name = nameFromPath(path)
	}

	return name, path
}

func nameFromPath(path string) string {
	name := filepath.Base(path)
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[:idx]
	}

	return removeIllegalChars(name)
}

func removeIllegalChars(name string) string {
	var result []rune
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			result = append(result, r)
		}
	}
	return string(result)
}
