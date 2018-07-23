package csvquery

import (
	"encoding/csv"
	"fmt"
	"io"

	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

type csvRowIter struct {
	closer io.Closer // free resources after the iter is done
	unlock func()    // unlock the resources after iter is done
	r      *csv.Reader
	closed bool
}

func (i *csvRowIter) Next() (sql.Row, error) {
	record, err := i.r.Read()
	if err != nil {
		_ = i.Close()
		if err == io.EOF {
			return nil, err
		}

		return nil, fmt.Errorf("csvquery: error reading record: %s", err)
	}

	var row = make(sql.Row, len(record))
	for i, v := range record {
		row[i] = v
	}

	return row, nil
}

func (i *csvRowIter) Close() error {
	if i.closed {
		return nil
	}

	i.closed = true
	i.unlock()
	return i.closer.Close()
}
