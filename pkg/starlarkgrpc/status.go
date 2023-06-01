package starlarkgrpc

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/grpc/codes"
)

var Status = starlarkstruct.FromStringDict(
	Symbol("grpc.status"),
	starlark.StringDict{
		"OK":                 starlark.MakeInt(int(codes.OK)),
		"CANCELED":           starlark.MakeInt(int(codes.Canceled)),
		"UNKNOWN":            starlark.MakeInt(int(codes.Unknown)),
		"INVALIDARGUMENT":    starlark.MakeInt(int(codes.InvalidArgument)),
		"DEADLINEEXCEEDED":   starlark.MakeInt(int(codes.DeadlineExceeded)),
		"NOTFOUND":           starlark.MakeInt(int(codes.NotFound)),
		"ALREADYEXISTS":      starlark.MakeInt(int(codes.AlreadyExists)),
		"PERMISSIONDENIED":   starlark.MakeInt(int(codes.PermissionDenied)),
		"RESOURCEEXHAUSTED":  starlark.MakeInt(int(codes.ResourceExhausted)),
		"FAILEDPRECONDITION": starlark.MakeInt(int(codes.FailedPrecondition)),
		"ABORTED":            starlark.MakeInt(int(codes.Aborted)),
		"OUTOFRANGE":         starlark.MakeInt(int(codes.OutOfRange)),
		"UNIMPLEMENTED":      starlark.MakeInt(int(codes.Unimplemented)),
		"INTERNAL":           starlark.MakeInt(int(codes.Internal)),
		"UNAVAILABLE":        starlark.MakeInt(int(codes.Unavailable)),
		"DATALOSS":           starlark.MakeInt(int(codes.DataLoss)),
		"UNAUTHENTICATED":    starlark.MakeInt(int(codes.Unauthenticated)),
	},
)
