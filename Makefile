.PHONY: server
server:
	bazel build \
		//example/routeguide:routeguide_proto_descriptor \
		//cmd/grpcstar


.PHONY: serve
serve: server
	GODEBUG=http2debug=2 \
	bazel-bin/cmd/grpcstar/grpcstar_/grpcstar \
		-protoset=bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		cmd/grpcstar/testdata/headers.grpc.star

.PHONY: routeguide_proto_descriptor
routeguide_proto_descriptor:
	bazel build //example/routeguide:routeguide_proto_descriptor
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/program/
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/starlarkgrpc/

.PHONY: mocks
mocks:
	mockery --srcpkg=google.golang.org/grpc --name=ClientStream

update_goldens:
	bazel run //cmd/grpcstar:grpcstar_test --action_env='GODEBUG=http2debug=1' -- --update
