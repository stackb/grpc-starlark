package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc/metadata"
)

var metadataSymbol = Symbol("Metadata")

type md struct {
	*starlarkstruct.Struct
	md metadata.MD
}

func newMetadata(meta metadata.MD) *md {
	return &md{
		md: meta,
		Struct: starlarkstruct.FromStringDict(
			metadataSymbol,
			starlark.StringDict{},
		),
	}
}

// Get implements part of the starlark.Mapping interface.
func (md *md) Get(k starlark.Value) (v starlark.Value, found bool, err error) {
	key := k.String()
	if vals, ok := md.md[key]; ok {
		if len(vals) == 1 {
			return starlark.String(vals[0]), true, nil
		}
		return goStringSliceToStarlarkList(vals), true, nil
	}
	return nil, false, nil
}

// SetKey implements part of the starlark.HasSetKey interface.
func (md *md) SetKey(k, v starlark.Value) (err error) {
	key := k.String()
	md.md.Set(key, v.String())
	return nil
}
