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
func (l *listener) String() string {
	return fmt.Sprintf("<net.Listener %v %v>", l.Listener.Addr().Network(), l.Listener.Addr())
}

// Type implements part of the starlark.Value interface
func (*listener) Type() string { return "net.Listener" }

// Freeze implements part of the starlark.Value interface
func (*listener) Freeze() {} // immutable

// Truth implements part of the starlark.Value interface
func (*listener) Truth() starlark.Bool { return starlark.False }

// Hash implements part of the starlark.Value interface
func (c *listener) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", c.Type())
}

// AttrNames implements part of the starlark.HasAttrs interface
func (c *listener) AttrNames() []string {
	return []string{"address", "network"}
}

// Attr implements part of the starlark.HasAttrs interface
func (c *listener) Attr(name string) (starlark.Value, error) {
	switch name {
	case "address":
		return starlark.String(c.Listener.Addr().String()), nil
	case "network":
		return starlark.String(c.Listener.Addr().Network()), nil
	}
	return nil, nil
}

func newListener(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var network string
	var address string
	if err := starlark.UnpackArgs("net.Listener", args, kwargs,
		"network?", &network,
		"address?", &address,
	); err != nil {
		return nil, err
	}
	if network == "" {
		network = "tcp"
	}
	if address == "" {
		address = "127.0.0.1:0"
	}

	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	value := &listener{
		Listener: l,
	}
	return value, nil
}
