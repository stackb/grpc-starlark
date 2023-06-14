BAZEL=bzl

.PHONY: build
build:
	$(BAZEL) build ...

.PHONY: test
test:
	$(BAZEL) test ... --runs_per_test=30

.PHONY: tidy
tidy:
	go mod tidy
	$(BAZEL) run update_go_repositories
	$(BAZEL) run gazelle

golden:
	$(BAZEL) run //cmd/grpcstar:grpcstar_test \
		--action_env=NOGODEBUG=http2debug=2 \
		-- \
		--update

.PHONY: serve
serve: build
	bazel-bin/cmd/grpcstar/grpcstar_/grpcstar \
		-p bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		-f cmd/grpcstar/testdata/routeguide.grpc.star

.PHONY: routeguide_proto_descriptor
routeguide_proto_descriptor:
	$(BAZEL) build //example/routeguide:routeguide_proto_descriptor
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/starlarkgrpc/
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/protodescriptorset/

.PHONY: plugin_proto_descriptor
plugin_proto_descriptor:
	$(BAZEL) build @protoapis//google/protobuf/compiler:plugin_descriptor
	cp -f bazel-bin/external/protoapis/google/protobuf/compiler/plugin_descriptor.pb cmd/protoc-gen-starlark
