package csvd

import (
	"fmt"
	"strings"
)

// CellError uses context to wrap cell value parse error.
type CellError struct {
	// The key used to get the cell value in the row.
	Key string
	// The value obtained from the cell by the given key in the row.
	// The value is empty string, if the cell cannot be found by the key.
	Val string

	// The row number of cell in the csv file.
	//
	// Numbering of rows starts at 1;
	Row int
	// The column number of cell in the csv file.
	//
	// Numbering of column starts at 1;
	// The column is -1, if the cell cannot be found by the key.
	Column int
	// The line number of cell in csv file.
	//
	// Numbering of lines starts at 1;
	// A row can span multiple lines.
	// The line is -1, if the cell cannot be found by the key.
	Line int

	// ParseErr occurs when parsing cell value.
	ParseErr error
}

func (e *CellError) writeTo(s *strings.Builder) {
	fmt.Fprintf(s, "row: %d, ", e.Row)
	fmt.Fprintf(s, "column: %d, ", e.Column)
	fmt.Fprintf(s, "key: %q, ", e.Key)
	fmt.Fprintf(s, "val: %q, ", e.Val)
	fmt.Fprintf(s, "err: %v", e.ParseErr)
}

func (e *CellError) Error() string {
	var s strings.Builder
	s.WriteString("cell err, ")
	e.writeTo(&s)
	return s.String()
}

func (e *CellError) Unwrap() error {
	return e.ParseErr
}

// Err returns a not nil error if both of e and ParseErr of e are not nil.
func (e *CellError) Err() error {
	if e == nil {
		return nil
	}
	if e.ParseErr == nil {
		return nil
	}
	return e
}

// RowError represents a set of cell value parse errors.
type RowError struct {
	CellErrors []*CellError
}

func (e *RowError) Error() string {
	var s strings.Builder
	fmt.Fprintf(&s, "%d cell error(s) found.\n", len(e.CellErrors))
	for _, cellError := range e.CellErrors {
		s.WriteString("- ")
		cellError.writeTo(&s)
		s.WriteString(".\n")
	}
	return s.String()
}

// Add adds cellErr to e if cellErr is not nil.
func (e *RowError) Add(cellErr *CellError) {
	if cellErr == nil {
		return
	}
	e.CellErrors = append(e.CellErrors, cellErr)
}

// Err returns a not nil error if e is not nil and length of CellErrors is not zero.
func (e *RowError) Err() error {
	if e == nil {
		return nil
	}
	if len(e.CellErrors) == 0 {
		return nil
	}
	return e
}
