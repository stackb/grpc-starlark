package starlarkgrpc

import (
	"context"
	"fmt"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc/metadata"
)

var contextSymbol = Symbol("Context")

type ctx struct {
	*starlarkstruct.Struct
	ctx context.Context
}

// Get implements part of the starlark.Mapping interface.  If the stored value
// is a starlark.Value it is returned directly, otherwise it will be converted
// to a starlark.String.
func (c *ctx) Get(key starlark.Value) (v starlark.Value, found bool, err error) {
	k := key.String()
	if val := c.ctx.Value(k); val != nil {
		value, ok := val.(starlark.Value)
		if !ok {
			value = starlark.String(fmt.Sprintf("%s", val))
		}
		return value, true, err
	}
	return nil, false, nil
}

func newCtx(x context.Context) *ctx {
	return &ctx{
		ctx: x,
		Struct: starlarkstruct.FromStringDict(
			contextSymbol,
			starlark.StringDict{
				"metadata": starlark.NewBuiltin(string(contextSymbol)+".metadata", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					if md, ok := metadata.FromIncomingContext(x); ok {
						return newMetadata(md), nil
					} else {
						return nil, nil
					}
				}),
				"with_timeout": starlark.NewBuiltin(string(contextSymbol)+".with_timeout", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var (
						millis int
					)
					if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
						"millis", &millis,
					); err != nil {
						return nil, err
					}

					next, cancel := context.WithTimeout(x, time.Millisecond*time.Duration(millis))

					return starlark.Tuple{
						newCtx(next),
						starlark.NewBuiltin(fn.Name()+".cancel", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
							cancel()
							return starlark.None, nil
						}),
					}, nil
				}),
				"with_value": starlark.NewBuiltin(string(contextSymbol)+".with_value", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
					var (
						key   string
						value starlark.Value
					)
					if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
						"key", &key,
						"value", &value,
					); err != nil {
						return nil, err
					}

					next := context.WithValue(x, key, value)

					return starlark.Tuple{newCtx(next)}, nil
				}),
			},
		),
	}
}
