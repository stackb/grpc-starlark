package starlarkgrpc

import (
	"fmt"

	"go.starlark.net/starlark"
)

type clientStreamCloseSend struct {
	closeSend func() error
}

func (cscs *clientStreamCloseSend) String() string {
	return fmt.Sprintf("<%s>", cscs.Type())
}

func (cscs *clientStreamCloseSend) Type() string { return fmt.Sprintf("%T", cscs) }

func (*clientStreamCloseSend) Truth() starlark.Bool { return starlark.True }

func (*clientStreamCloseSend) Freeze() {} // immutable

func (csr *clientStreamCloseSend) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", csr.Type())
}

func (css *clientStreamCloseSend) Name() string {
	return "ClientStreamCloseSend"
}

func (cscs *clientStreamCloseSend) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := cscs.closeSend(); err != nil {
		return nil, err
	}
	return starlark.None, nil
}
