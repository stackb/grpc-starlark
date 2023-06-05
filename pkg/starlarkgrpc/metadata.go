package starlarkgrpc

import (
	"fmt"
	"sort"

	"go.starlark.net/starlark"
	"google.golang.org/grpc/metadata"
)

type md metadata.MD

func metadataFromValue(v starlark.Value) (md, bool) {
	if meta, ok := v.(md); ok {
		return meta, true
	}
	if dict, ok := v.(*starlark.Dict); ok {
		return metadataFromDict(dict), true
	}
	return nil, false
}

func metadataFromDict(dict *starlark.Dict) md {
	meta := metadata.Pairs()
	for _, k := range dict.Keys() {

		var key string
		switch t := k.(type) {
		case starlark.String:
			key = t.GoString()
		default:
			key = t.String()
		}

		v, _, _ := dict.Get(k)
		var val string
		switch t := v.(type) {
		case starlark.String:
			val = t.GoString()
		default:
			val = t.String()
		}

		meta[key] = []string{val}
	}
	return md(meta)
}

func newMetadata(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return makeMetadata(metadata.New(map[string]string{})), nil
}

func makeMetadata(meta metadata.MD) md {
	return md(meta)
}

// String implements the Stringer interface.
func (d md) String() string { return fmt.Sprintf("md<%v>", d.AttrNames()) }

// Type returns a short string describing the value's type.
func (d md) Type() string { return "metadata.MD" }

// Freeze renders *md immutable. required by starlark.Value interface
// because duration is already immutable this is a no-op.
func (d md) Freeze() {}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y)
// required by starlark.Value interface.
func (d md) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", d.Type())
}

// Truth reports whether the duration is non-zero.
func (d md) Truth() starlark.Bool { return starlark.True }

// AttrNames lists available dot expression strings. required by
// starlark.HasAttrs interface.
func (d md) AttrNames() (names []string) {
	for k := range d {
		names = append(names, k)
	}
	sort.Strings(names)
	return
}

// Attr gets a value for a string attribute, implementing dot expression support
// in starklark. required by starlark.HasAttrs interface.
func (d md) Attr(name string) (starlark.Value, error) {
	if vals, ok := d[name]; ok {
		if len(vals) == 1 {
			return starlark.String(vals[0]), nil
		}
		return goStringSliceToStarlarkList(vals), nil

	}
	return nil, fmt.Errorf("unrecognized %s attribute %q", d.Type(), name)
}

// Get implements part of the starlark.Mapping interface.
func (d md) Get(k starlark.Value) (v starlark.Value, found bool, err error) {
	key := k.String()
	if vals, ok := d[key]; ok {
		if len(vals) == 1 {
			return starlark.String(vals[0]), true, nil
		}
		return goStringSliceToStarlarkList(vals), true, nil
	}
	return nil, false, nil
}

// SetKey implements part of the starlark.HasSetKey interface.
func (d md) SetKey(k, v starlark.Value) (err error) {
	var key, val string

	switch t := k.(type) {
	case starlark.String:
		key = t.GoString()
	default:
		key = v.String()
	}

	switch t := v.(type) {
	case starlark.String:
		val = t.GoString()
	default:
		val = v.String()
	}

	d[key] = []string{val}

	return nil
}
