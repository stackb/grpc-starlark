.PHONY: build
build:
	bazel build ...

.PHONY: test
test:
	bazel test ... --runs_per_test=30

golden:
	bazel run //cmd/grpcstar:grpcstar_test \
		--action_env='GODEBUG=http2debug=1' \
		-- \
		--update

.PHONY: serve
serve: build
	GODEBUG=http2debug=2 \
	bazel-bin/cmd/grpcstar/grpcstar_/grpcstar \
		-p bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		-f cmd/grpcstar/testdata/routeguide.grpc.star

.PHONY: routeguide_proto_descriptor
routeguide_proto_descriptor:
	bazel build //example/routeguide:routeguide_proto_descriptor
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/starlarkgrpc/
