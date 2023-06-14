# rules

The `@build_stack_grpc_starlark//rules:` package contains custom bazel rules for
working with grpc-starlark.

## `grpcstar_binary`

The `grpcstar_binary` generates a standalone binary with the descriptor and
entrypoint script embedded into the binary.  This packages the required files
into an easily runnable / deployable executable that does not require additional
command line arguments.

Example:

```py
load("@build_stack_grpc_starlark//rules:grpcstar_binary.bzl", "grpcstar_binary")

grpcstar_binary(
    name = "server",
    descriptor = ":routeguide_proto",
    main = "routeguide.main.star",
)
```

### Attributes

| name         | type   | required | desciption                                                  |
| ------------ | ------ | -------- | ----------------------------------------------------------- |
| `main`       | label  | yes      | The starlark source file having a `main(ctx)` func          |
| `descriptor` | label  | yes      | The proto_descriptor_set file                               |
| `template`   | label  | no       | The template file that produces a `main.go`                 |
| `importpath` | string | no       | for `go_library.importpath`, defaults to `{package}/{name}` |

### Files 

> This information is not necessary for using the rule but helps explain how it works.

```sh
$ bazel query '//example/routeguide:*' --output label_kind
source file //example/routeguide:routeguide.proto          # source file: input for proto_library.srcs
proto_library rule //example/routeguide:routeguide_proto   # generates: the compiled descriptor.pb, for genrule
genrule rule //example/routeguide:server_descriptor        # generates: copy of descriptor file (same package, easier for embedding)
generated file //example/routeguide:server.descriptor_     # generated file: input for go_library.embedsrcs
source file //example/routeguide:routeguide.main.star      # source file: input for genrule
genrule rule //example/routeguide:server_star              # generates: copy of source file (same package, easier for embedding)
generated file //example/routeguide:server.star_           # generated file: input for go_library.embedsrcs
_grpcstar_entrypoint rule //example/routeguide:server_main # generates: main.go file, for go_library.srcs
generated file //example/routeguide:server_main.go         # generated file: input for go_library.srcs
go_library rule //example/routeguide:server_lib            # generates: archive for go_binary.embed
go_binary rule //example/routeguide:server                 # generates: the executable, for running!
```
