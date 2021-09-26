package csvd

// Row represents a row in csv file.
type Row struct {
	d      *Decoder
	record []string

	rowError       *RowError
	lastParseError *CellError
}

func (r *Row) check() {
	d := r.d
	if d == nil {
		panic("row without decoder, please get Row from the Decoder")
	}
	if r.d.header == nil {
		panic("row without header, forgot to call ParseHeader()?")
	}
}

// Record returns the original data in the row.
func (r *Row) Record() []string {
	return r.record
}

func (r *Row) get(key string) (int, string) {
	r.check()

	idx := r.d.header.Get(key)
	var val string
	if idx != -1 {
		val = r.record[idx]
	}
	return idx, val
}

// Get gets the value of cell associated with the given key.
// If there are no cell associated with the key, Get returns empty string.
//
// The key is canonicalized by CanonicalHeaderKey.
func (r *Row) Get(key string) string {
	_, val := r.get(key)
	return val
}

// Has checks if has a cell associated with the given key.
// If there are no cell associated with the key, Has returns false.
//
// The key is canonicalized by CanonicalHeaderKey.
func (r *Row) Has(key string) bool {
	r.check()

	return r.d.header.Has(key)
}

// LastParseError returns the error that occurred the last time the Parse method was called.
func (r *Row) LastParseError() *CellError {
	if r.lastParseError != nil {
		return r.lastParseError
	}
	return nil
}

func (r *Row) addCellError(cellErr *CellError) {
	if r.rowError == nil {
		r.rowError = &RowError{}
	}
	r.rowError.Add(cellErr)
}

// Error returns all errors that occurred when Parse was called.
func (r *Row) Error() *RowError {
	return r.rowError
}

// Parse gets the value of cell associated with the given key then
// call parse method with that value.
//
// The error returned by the parse method is wrapped as a CellError and
// can be obtained through the LastParseError method.
//
// If you want to get all parse errors, call the Error() method.
func (r *Row) Parse(key string, parse func(val string) error) {
	r.check()

	idx, val := r.get(key)
	var lastParseError *CellError
	err := parse(val)
	if err != nil {
		line, _ := r.d.reader.FieldPos(idx)
		lastParseError = &CellError{
			Key:      key,
			Val:      val,
			Row:      r.d.rowNumber,
			Column:   idx + 1,
			Line:     line,
			ParseErr: err,
		}
	}

	r.addCellError(lastParseError)
	r.lastParseError = lastParseError
}
