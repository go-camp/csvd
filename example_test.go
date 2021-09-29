package csvd_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-camp/csvd"
)

func ExampleDecoder() {
	var csvFile = "User  \tName , Age\ninvalid int,a\ngood age , 18\ntoo yong,17\ntoo old,121\n"
	r := bytes.NewReader([]byte(csvFile))
	decoder, err := csvd.NewDecoder(csvd.Options{
		Reader: csv.NewReader(r),
	})
	if err != nil {
		panic(err)
	}
	decoder.ParseHeader()
	for {
		row, err := decoder.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		row.Parse("user name", func(val string) error {
			return nil
		})
		row.Parse("age", func(val string) error {
			val = strings.TrimSpace(val)

			age, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				return err
			}
			if age < 18 {
				return errors.New("too yong")
			}
			if age > 120 {
				return errors.New("too old")
			}
			return nil
		})
		if err := row.Error().Err(); err != nil {
			fmt.Println(err)
		}
	}
	// output:
	// 1 cell error(s) found.
	// - row: 2, column: 2, key: "age", val: "a", err: strconv.ParseInt: parsing "a": invalid syntax.
	//
	// 1 cell error(s) found.
	// - row: 4, column: 2, key: "age", val: "17", err: too yong.
	//
	// 1 cell error(s) found.
	// - row: 5, column: 2, key: "age", val: "121", err: too old.
}
