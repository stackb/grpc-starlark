package starlarkgrpc

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"
)

func newPerRpcCredentials(parent *starlark.Thread, call starlark.Callable) *perRpcCredentials {
	return &perRpcCredentials{
		parent:                    parent,
		call:                      call,
		requiresTransportSecurity: false, // TODO
	}
}

type perRpcCredentials struct {
	parent                    *starlark.Thread
	call                      starlark.Callable
	requiresTransportSecurity bool
}

func (r *perRpcCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	args := starlark.Tuple{
		newCtx(ctx),
		goStringSliceToStarlarkList(uri),
	}
	value, err := starlark.Call(newThread(r.parent, r.call.String()), r.call, args, []starlark.Tuple{})
	if err != nil {
		return nil, fmt.Errorf("perRpcCredentials error %s: %w", r.call.String(), err)
	}
	dict, ok := value.(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("perRpcCredentials return value must be a dict, got %s (%T)", value, value)
	}
	return starlarkDictToGoStringMap(dict), nil
}

func (r *perRpcCredentials) RequireTransportSecurity() bool {
	return r.requiresTransportSecurity
}

func goStringSliceToStarlarkList(list []string) *starlark.List {
	elems := make([]starlark.Value, len(list))
	for i, v := range list {
		elems[i] = starlark.String(v)
	}
	return starlark.NewList(elems)
}
func starlarkDictToGoStringMap(dict *starlark.Dict) map[string]string {
	md := make(map[string]string)
	for _, k := range dict.AttrNames() {
		if v, err := dict.Attr(k); err != nil {
			md[k] = v.String()
		}
	}
	return md
}
