package starlarkos

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.starlark.net/starlark"
)

// File is the starlark.Value representation of a file.
type File struct {
	*os.File
	data []byte
}

func (f File) String() string        { return string(f.data) }
func (f File) GoString() string      { return string(f.data) }
func (f File) Type() string          { return "os.file" }
func (f File) Freeze()               {} // immutable
func (f File) Truth() starlark.Bool  { return true }
func (f File) Hash() (uint32, error) { return starlark.String(f.String()).Hash() }
func (f File) Len() int              { return len(f.data) } // bytes

func getStdin(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs("os.stdin", args, kwargs); err != nil {
		return nil, err
	}
	if data, err := ioutil.ReadAll(os.Stdin); err != nil {
		return nil, fmt.Errorf("reading stdin: %w", err)
	} else {
		return &File{File: os.Stdin, data: data}, nil
	}
}
