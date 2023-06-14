# rules

The `@build_stack_grpc_starlark//rules:` package contains custom bazel rules for
working with grpc-starlark.

## `WORKSPACE`

In order to consume grpc-starlark in a bazel workspace, use something like:

```py
# Release: v0.6.0
# TargetCommitish: master
# Date: 2023-06-14 14:56:20 +0000 UTC
# URL: https://github.com/stackb/grpc-starlark/releases/tag/v0.6.0
# Size: 89768 (90 kB)
http_archive(
    name = "build_stack_grpc_starlark",
    sha256 = "e0e4310c4b968277f68f99206d38b0fb3c3aff36fae8a8a8daab9d422d88dc50",
    strip_prefix = "grpc-starlark-0.6.0",
    urls = ["https://github.com/stackb/grpc-starlark/archive/v0.6.0.tar.gz"],
)
```

If you need the go dependencies, use something like:

```py
load("@build_stack_grpc_starlark//:go_repositories.bzl", build_stack_grpc_starlark_go_repositories = "go_repositories")

build_stack_grpc_starlark_go_repositories()
```

> Not required if you have a workspace that is using go imports from
> grpc-starlark and they are already in your `go.mod` file, and you have a
> workflow like `gazelle update-repos`.

If you need base dependencies (rules_go, etc), use something like:

```py
load("@build_stack_grpc_starlark//:repositories.bzl", build_stack_grpc_starlark_repositories = "repositories")

build_stack_grpc_starlark_repositories()
```

> Not required if you already have a workspace using rules_go.

A note about c++ toolchains.  If you are running more recent versions of rules_go and bazel and are seeing:

```
external/io_bazel_rules_go/BUILD.bazel:86:17: While resolving toolchains for target @io_bazel_rules_go//:cgo_context_data: No matching toolchains found for types @bazel_tools//tools/cpp:toolchain_type.
```

This is discussed more in https://github.com/bazelbuild/rules_go/issues/3470.

As a workaround, you might be able provide a set of C++ toolchains using zig:

```py
# ----------------------------------------------------
# @hermetic_cc_toolchain (zig)
# ----------------------------------------------------

load("@hermetic_cc_toolchain//toolchain:defs.bzl", zig_toolchains = "toolchains")

# Plain zig_toolchains() will pick reasonable defaults. See
# toolchain/defs.bzl:toolchains on how to change the Zig SDK version and
# download URL.
zig_toolchains()

register_toolchains(
    "@zig_sdk//toolchain:linux_amd64_gnu.2.28",
    "@zig_sdk//toolchain:linux_arm64_gnu.2.28",
    "@zig_sdk//toolchain:darwin_amd64",
    "@zig_sdk//toolchain:darwin_arm64",
    "@zig_sdk//toolchain:windows_amd64",
    "@zig_sdk//toolchain:windows_arm64",
)
```

> @hermetic_cc_toolchain is declared in `repositories.bzl`.  Additional flags in `.bazelrc` will be needed.

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

```
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

## `grpcstar_image`

The `grpcstar_image` rule generates a container image.  It has usage similar to `grpcstar_binary`, although it does not use that rule directly.

Example:

```py
load("@build_stack_grpc_starlark//rules:grpcstar_image.bzl", "grpcstar_image")

grpcstar_image(
    name = "image",
    descriptor = ":routeguide_proto",
    main = "routeguide.main.star",
)
```

### Attributes

| name         | type  | required | desciption                                                                  |
| ------------ | ----- | -------- | --------------------------------------------------------------------------- |
| `main`       | label | yes      | The starlark source file having a `main(ctx)` func                          |
| `descriptor` | label | yes      | The proto_descriptor_set file                                               |
| `executable` | label | no       | The grpcstar binary, defaults to `@build_stack_grpc_starlark//cmd/grpcstar` |
| `base`       | label | no       | The base image to use, defaults to `@go_image_base//image`                  |

The container layers will be:

1. the base layer
2. the executable
3. the descriptor
4. entrypoint script

You are most likely to update only theqentrypoint script during development,
making container pushes and pulls fast.