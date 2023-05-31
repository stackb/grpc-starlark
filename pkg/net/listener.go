package net

import (
	"fmt"
	"net"

	"go.starlark.net/starlark"
)

// listener implements starlark.Value for a net.Listener.
type listener struct {
	net.Listener
}

// String implements part of the starlark.Value interface
func (*listener) String() string { return "listener" }

// Type implements part of the starlark.Value interface
func (*listener) Type() string { return "listener" }

// Freeze implements part of the starlark.Value interface
func (*listener) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*listener) Truth() starlark.Bool { return starlark.False }

// Hash implements part of the starlark.Value interface
func (c *listener) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

func newListener(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

	l, err := net.Listen(network, address)
	if err != nil {
		return nil, fmt.Errorf("starting listener on %s: %w", address, err)
	}
	value := &listener{
		Listener: l,
	}
	return value, nil
}
