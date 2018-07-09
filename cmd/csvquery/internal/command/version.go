package command

import "fmt"

// Version is a command that will output the version.
type Version struct {
	Version, Commit, Date string
}

// Execute the command.
func (v Version) Execute([]string) error {
	fmt.Printf("%v, commit %v, built at %v\n", v.Version, v.Commit, v.Date)
	return nil
}
