package starlarkgrpc

import "go.starlark.net/starlark"

// Symbol is the type of a Starlark constructor symbol.  It prints more
// favorably than a starlark.String.
type Symbol string

func (s Symbol) String() string             { return string(s) }
func (s Symbol) GoString() string           { return string(s) }
func (s Symbol) Type() string               { return "symbol" }
func (s Symbol) Freeze()                    {} // immutable
func (s Symbol) Truth() starlark.Bool       { return len(s) > 0 }
func (s Symbol) Hash() (uint32, error)      { return starlark.String(s).Hash() }
func (s Symbol) Len() int                   { return len(s) } // bytes
func (s Symbol) Index(i int) starlark.Value { return s[i : i+1] }
