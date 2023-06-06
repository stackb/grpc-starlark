package program

import (
	_ "embed"
)

//go:embed routeguide_proto_descriptor.pb
var routeguideProtoDescriptor []byte

//go:embed routeguide.grpc.star
var routeguideGrpcStar string
