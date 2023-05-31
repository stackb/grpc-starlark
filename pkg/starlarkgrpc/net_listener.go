package starlarkgrpc

import (
	"fmt"
	"net"

	"go.starlark.net/starlark"
)

// netListener implements starlark.Value for a net.Listener.
type netListener struct {
	net.Listener
}

// String implements part of the starlark.Value interface
func (*netListener) String() string { return "netListener" }

// Type implements part of the starlark.Value interface
func (*netListener) Type() string { return "netListener" }

// Freeze implements part of the starlark.Value interface
func (*netListener) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*netListener) Truth() starlark.Bool { return starlark.False }

// Hash implements part of the starlark.Value interface
func (c *netListener) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func newNetListenerFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var network string
		var address string
		if err := starlark.UnpackArgs("Listener", args, kwargs,
			"network?", &network,
			"address?", &address,
		); err != nil {
			return nil, err
		}
		if network == "" {
			network = "tcp"
		}
		if address == "" {
			address = "127.0.0.0:0"
		}

		listener, err := net.Listen(network, address)
		if err != nil {
			return nil, fmt.Errorf("starting listener on %s: %w", address, err)
		}
		value := &netListener{
			Listener: listener,
		}
		return value, nil
	}
}
