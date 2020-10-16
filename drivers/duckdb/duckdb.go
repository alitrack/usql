package duckdb

import (
	"strings"

	// DRIVER: duckdb
	_ "github.com/marcboeker/go-duckdb"
	"github.com/xo/usql/drivers"
	"github.com/xo/xoutil"
)

func init() {
	drivers.Register("duckdb", drivers.Driver{
		Name:                   "duckdb",
		AllowMultilineComments: true,
		LexerName:              "duckdb",
		ForceParams: drivers.ForceQueryParameters([]string{
			"loc", "auto",
		}),
		Version: func(db drivers.DB) (string, error) {
			var ver string
			err := db.QueryRow(`select library_version from pragma_version();`).Scan(&ver)
			if err != nil {
				return "", err
			}
			return "duckdb " + ver, nil
		},
		Err: func(err error) (string, string) {

			// if e, ok := err.(sqlite3.Error); ok {
			// 	return strconv.Itoa(int(e.Code)), e.Error()
			// }

			code, msg := "", err.Error()
			// if e, ok := err.(sqlite3.ErrNo); ok {
			// 	code = strconv.Itoa(int(e))
			// }

			return code, msg
		},
		ConvertBytes: func(buf []byte, tfmt string) (string, error) {
			// attempt to convert buf if it matches a time format, and if it
			// does, then return a formatted time string.
			s := string(buf)
			if s != "" && strings.TrimSpace(s) != "" {
				t := new(xoutil.SqTime)
				if err := t.Scan(buf); err == nil {
					return t.Format(tfmt), nil
				}
			}
			return s, nil
		},
	})
	// duckdbSchema := dburl.Scheme{"duckdb", dburl.GenOpaque, 0, true, []string{"duckdb", "file"}, ""}
	// dburl.Register(duckdbSchema)
}
