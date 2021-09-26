package csvd

import (
	"encoding/csv"
	"errors"
)

// Decoder reads rows from a CSV-encoded file.
type Decoder struct {
	reader *csv.Reader

	header       Header
	headerRecord []string

	rowNumber int
}

type Options struct {
	Reader *csv.Reader
}

// NewDecoder returns a new Decoder that reads from opts.Reader.
func NewDecoder(opts Options) (*Decoder, error) {
	if opts.Reader == nil {
		return nil, errors.New("Options.Reader can't be nil")
	}
	return &Decoder{
		reader: opts.Reader,
	}, nil
}

// ParseHeader parses a row of csv file into header.
func (d *Decoder) ParseHeader() error {
	record, err := d.reader.Read()
	if err != nil {
		return err
	}

	d.rowNumber++

	d.header = ParseHeader(record)
	d.headerRecord = append([]string(nil), record...)
	return nil
}

// Next reads one row of csv file.
//
// If there is no row left to be read, Next returns nil, io.EOF.
func (d *Decoder) Next() (*Row, error) {
	record, err := d.reader.Read()
	if err != nil {
		return nil, err
	}

	d.rowNumber++

	return &Row{
		d:      d,
		record: record,
	}, nil
}
