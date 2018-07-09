package command

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testData(file string) string {
	return filepath.Join("..", "..", "..", "..", "_testdata", file)
}

func TestREPL(t *testing.T) {
	testCases := []struct {
		in, out string
	}{
		{"help", replHelp + "\n"},
		{"load " + testData("ratings.csv"), "Table loaded successfully.\n"},
		{
			"load " + testData("foo.csv"),
			"unable to add table \"foo\": open ../../../../_testdata/foo.csv: no such file or directory\n",
		},
		{
			"SELECT * FROM superheroes",
			"" +
				"+----+--------------+---------+\n" +
				"| ID |     NAME     | COMPANY |\n" +
				"+----+--------------+---------+\n" +
				"| sm | Superman     | dc      |\n" +
				"| bm | Batman       | dc      |\n" +
				"| ww | Wonder Woman | dc      |\n" +
				"| ma | Miss America | marvel  |\n" +
				"+----+--------------+---------+\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			require := require.New(t)
			var in = bytes.NewBuffer([]byte(tt.in))
			var out bytes.Buffer
			err := (&REPL{
				Stdin:   in,
				Stdout:  &out,
				Stderr:  &out,
				baseCmd: baseCmd{Files: []string{testData("superheroes.csv")}},
			}).Execute(nil)
			require.NoError(err)
			require.Contains(out.String(), initMsg+"\n"+tt.out)
		})
	}
}
