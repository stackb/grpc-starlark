.PHONY: server
server:
	bazel build \
		//example/routeguide:routeguide_proto_descriptor \
		//cmd/grpcstar


.PHONY: serve
serve:
	bazel-bin/cmd/grpcstar/grpcstar_/grpcstar \
		-protoset=bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		pkg/program/routeguide.grpc.star

.PHONY: routeguide_proto_descriptor
routeguide_proto_descriptor:
	bazel build //example/routeguide:routeguide_proto_descriptor
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/program/

.PHONY: mocks
mocks:
	mockery --srcpkg=google.golang.org/grpc --name=ClientStream
