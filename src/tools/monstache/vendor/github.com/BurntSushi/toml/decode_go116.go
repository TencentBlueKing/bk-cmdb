// +build go1.16

package toml

import (
	"io/fs"
)

// DecodeFS is just like Decode, except it will automatically read the contents
// of the file at `fpath` from a fs.FS instance.
func DecodeFS(fsys fs.FS, fpath string, v interface{}) (MetaData, error) {
	fp, err := fsys.Open(fpath)
	if err != nil {
		return MetaData{}, err
	}
	defer fp.Close()
	return NewDecoder(fp).Decode(v)
}
