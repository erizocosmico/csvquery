package main

import (
	"os"

	"github.com/erizocosmico/csvquery/cmd/csvquery/internal/command"

	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	parser := flags.NewNamedParser("csvquery", flags.Default)

	_, err := parser.AddCommand(
		"server",
		"Start a MySQL-compatible server to query CSV files.",
		"",
		new(command.Server),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = parser.AddCommand(
		"repl",
		"Start a REPL to query CSV files.",
		"",
		&command.REPL{Stdin: os.Stdin, Stderr: os.Stderr, Stdout: os.Stdout},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = parser.AddCommand(
		"version",
		"Show version of the program.",
		"",
		&command.Version{Version: version, Commit: commit, Date: date},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}
