package command

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/erizocosmico/csvquery"
	"github.com/olekukonko/tablewriter"
	sqle "gopkg.in/src-d/go-mysql-server.v0"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

// REPL will start a REPL to query CSV files.
type REPL struct {
	baseCmd

	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

const initMsg = `
  ____ ______     _____                        
 / ___/ ___\ \   / / _ \ _   _  ___ _ __ _   _ 
| |   \___ \\ \ / / | | | | | |/ _ \ '__| | | |
| |___ ___) |\ V /| |_| | |_| |  __/ |  | |_| |
 \____|____/  \_/  \__\_\\__,_|\___|_|   \__, |
                                         |___/ 

Query CSV files with SQL.

Enter "help" to see more details about the usage, or "quit" to exit.
`

const replHelp = `Help:

  help                  Show this message.
  quit,exit             Exit the REPL.
  load <path>[ <name>]  Load a new CSV file as a table. If name is not
                        provided, the file name will be used as table
                        name removing all characters that are not alpha
                        numeric or underscores and removing the extension.
  <SQL query>           Execute a SQL query.
`

var errHelp = errors.New(replHelp)
var loadRegex = regexp.MustCompile(`^\s?load\s+([^\s]+)\s?(.*)`)

// Execute the command
func (c *REPL) Execute([]string) error {
	engine, db, err := c.engine()
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "csvquery> ",
		Stdin:  c.Stdin,
		Stderr: c.Stderr,
		Stdout: c.Stdout,
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	fmt.Fprintln(c.Stdout, initMsg)

	for {
		line, err := rl.Readline()
		if err != nil {
			return nil
		}

		switch strings.ToLower(strings.TrimSpace(line)) {
		case "exit", "quit":
			fmt.Fprintln(c.Stdout, "Bye, have a nice day!")
			return nil
		case "help":
			fmt.Fprintln(c.Stdout, replHelp)
		default:
			if strings.HasPrefix(strings.TrimSpace(strings.ToLower(line)), "load ") {
				if err := loadFile(line, db); err != nil {
					fmt.Fprintln(c.Stderr, err)
				} else {
					fmt.Fprintln(c.Stdout, "Table loaded successfully.")
				}

				continue
			}

			c.runQuery(engine, line)
		}
	}
}

func loadFile(line string, db *csvquery.Database) error {
	if loadRegex.MatchString(line) {
		matches := loadRegex.FindStringSubmatch(line)
		var path, name string
		switch len(matches) {
		case 2:
			path = matches[1]
		case 3:
			path = matches[1]
			name = matches[2]
		default:
			return errHelp
		}

		if name == "" {
			name = nameFromPath(path)
		}

		return db.AddTable(name, path)
	}

	return errHelp
}

func (c *REPL) runQuery(engine *sqle.Engine, query string) {
	start := time.Now()
	schema, iter, err := engine.Query(sql.NewEmptyContext(), query)
	if err != nil {
		fmt.Fprintln(c.Stderr, err)
		return
	}

	var columns = make([]string, len(schema))
	for i, col := range schema {
		columns[i] = col.Name
	}

	writer := tablewriter.NewWriter(c.Stdout)
	writer.SetHeader(columns)

	var errors []string
	var rows int
	for {
		row, err := iter.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			errors = append(errors, err.Error())
			break
		}

		var values = make([]string, len(row))
		for i, v := range row {
			values[i] = fmt.Sprint(v)
		}

		rows++
		writer.Append(values)
	}

	if err := iter.Close(); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		fmt.Fprintln(c.Stderr, "Error:")
		fmt.Fprintf(c.Stderr, "%s\n", strings.Join(errors, "\n"))
	} else {
		writer.Render()
		fmt.Fprintf(c.Stdout, "%d row(s) in %v\n", rows, time.Since(start))
	}
}
