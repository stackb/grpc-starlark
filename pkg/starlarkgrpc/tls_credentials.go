package starlarkgrpc

import (
	"github.com/stackb/grpc-starlark/pkg/starlarkcrypto"
	"github.com/stackb/grpc-starlark/pkg/starlarkutil"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type transportCredentials struct {
	credentials.TransportCredentials
	*starlarkstruct.Struct
}

func newTlsCredentials(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var config starlarkcrypto.Config

	if err := starlark.UnpackArgs(fn.Name(), args, kwargs,
		"config", &config,
	); err != nil {
		return nil, err
	}

	return transportCredentials{
		TransportCredentials: credentials.NewTLS(config.Config),
		Struct: starlarkstruct.FromStringDict(
			starlarkutil.Symbol(fn.Name()),
			starlark.StringDict{
				"config": config,
			},
		),
	}, nil
}

func newInsecureCredentials(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs); err != nil {
		return nil, err
	}
	return transportCredentials{
		TransportCredentials: insecure.NewCredentials(),
		Struct: starlarkstruct.FromStringDict(
			starlarkutil.Symbol(fn.Name()),
			starlark.StringDict{},
		),
	}, nil
}
