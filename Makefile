descriptor_set:

grpc-starlark:

server:
	bazel build \
		//example/routeguide:routeguide_proto_descriptor \
		//cmd/grpc-starlark

serve:
	bazel-bin/cmd/grpc-starlark/grpc-starlark_/grpc-starlark \
		-protoset=bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		-load=example/routeguide/routeguide.grpc.star

client:
	cd example/module && npx tsc

perf:
	node example/module/dist/main.js