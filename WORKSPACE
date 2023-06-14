workspace(name = "build_stack_grpc_starlark")

load("//:repositories.bzl", "repositories")

repositories()

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

# ----------------------------------------------------
# @rules_proto
# ----------------------------------------------------

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

# ----------------------------------------------------
# @io_bazel_rules_go
# ----------------------------------------------------

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(version = "1.18.2")

# ----------------------------------------------------
# @build_stack_rules_proto
# ----------------------------------------------------

register_toolchains("@build_stack_rules_proto//toolchain:standard")

load("//:proto_repositories.bzl", "proto_repositories")

proto_repositories()

# ----------------------------------------------------
# @build_stack_rules_proto
# ----------------------------------------------------

load("@build_stack_rules_proto//:go_deps.bzl", "gazelle_protobuf_extension_go_deps")

gazelle_protobuf_extension_go_deps()

load("@build_stack_rules_proto//deps:go_core_deps.bzl", "go_core_deps")

go_core_deps()

# ----------------------------------------------------
# external go dependencies
# ----------------------------------------------------

load("//:go_repositories.bzl", "go_repositories")

go_repositories()

# ----------------------------------------------------
# @bazel_gazelle
# ----------------------------------------------------
# gazelle:repository_macro go_repositories.bzl%go_repositories

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# ----------------------------------------------------
# @io_bazel_rules_docker
# ----------------------------------------------------

load("@io_bazel_rules_docker//go:image.bzl", _go_image_repos = "repositories")
load("@io_bazel_rules_docker//repositories:repositories.bzl", container_repositories = "repositories")

container_repositories()

_go_image_repos()
