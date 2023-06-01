package starlarkgrpc

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeGrpcError(strct *starlarkstruct.Struct) error {
	sym, ok := strct.Constructor().(Symbol)
	if !ok {
		return fmt.Errorf("unexpected handler return Struct constructor: want 'grpc.Error', got: %s", strct.Constructor())
	}
	if sym.String() != "Error" {
		return fmt.Errorf("unexpected handler return Struct: want 'grpc.Error', got: %s", strct.Constructor())
	}
	code, err := getStructAttrUint64(strct, "code")
	if err != nil {
		return fmt.Errorf("grpc.Error.code: %w", err)
	}
	message, err := getStructAttrString(strct, "message")
	if err != nil {
		return fmt.Errorf("grpc.Error.message: %w", err)
	}
	return status.Errorf(codes.Code(code), message)
}

func newError(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	code := int(codes.Unknown)
	var message string

	if err := starlark.UnpackArgs("grpc.Error", args, kwargs,
		"code?", &code,
		"message?", &message,
	); err != nil {
		return nil, err
	}
	value := starlarkstruct.FromStringDict(
		Symbol("grpc.Error"),
		starlark.StringDict{
			"code":    starlark.MakeInt(code),
			"message": starlark.String(message),
		},
	)

	return value, nil
}
