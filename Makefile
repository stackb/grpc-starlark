descriptor_set:

grpc-starlark:

server:
	bazel build \
		//example/routeguide:routeguide_proto_descriptor \
		//cmd/grpcstar

serve:
	bazel-bin/cmd/grpcstar/grpcstar_/grpcstar \
		-protoset=bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		pkg/program/routeguide.grpc.star

routeguide_proto_descriptor:
	bazel build //example/routeguide:routeguide_proto_descriptor
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/program/
