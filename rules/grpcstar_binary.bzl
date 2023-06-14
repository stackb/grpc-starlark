load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# gazelle:ignore testdata

def _grpcstar_entrypoint_impl(ctx):
    ctx.actions.expand_template(
        template = ctx.file.template,
        output = ctx.outputs.go,
        substitutions = {
            "{DESCRIPTOR_PATH}": ctx.file.descriptor.basename,
            "{MAIN_PATH}": ctx.file.main.basename,
        },
    )
    return [DefaultInfo(
        files = depset([ctx.outputs.go]),
    )]

_grpcstar_entrypoint = rule(
    implementation = _grpcstar_entrypoint_impl,
    attrs = {
        "descriptor": attr.label(
            doc = "a label pointing to a proto_library or proto_descriptor_set rule",
            allow_single_file = True,
            mandatory = True,
        ),
        "main": attr.label(
            doc = "a starlark source that contains the main entrypoint",
            allow_single_file = True,
            mandatory = True,
        ),
        "srcs": attr.label_list(
            doc = "a list of additional starlark source files that can be loaded by the main entrypoint",
            allow_files = True,
        ),
        "template": attr.label(
            allow_single_file = True,
            default = "@build_stack_grpc_starlark//rules:grpcstar_binary.go.tmpl",
        ),
    },
    outputs = {
        "go": "%{name}.go",
    },
)

def grpcstar_binary(**kwargs):
    """grpcstar_binary is a macro that generates a main.go source file and compiles a go_binary for it.

    Args:
        **kwargs: the keyword args dict
    Returns:
        None
    """
    name = kwargs.pop("name")
    goname = name + "_main"
    libname = name + "_lib"
    starname = name + "_star"
    dname = name + "_descriptor"

    scripts = kwargs.pop("scripts", [])
    srcs = kwargs.pop("srcs", [])
    deps = kwargs.pop("deps", [])
    visibility = kwargs.pop("visibility", ["//visibility:public"])

    importpath = kwargs.pop("importpath", "")
    if not importpath:
        fail("grpcstar_binary.importpath is required")
    descriptor = kwargs.pop("descriptor", "")
    if not descriptor:
        fail("grpcstar_binary.descriptor is required")
    main = kwargs.pop("main", "")
    if not main:
        fail("grpcstar_binary.main is required")

    native.genrule(
        name = starname,
        srcs = [main],
        outs = [name + ".star_"],
        cmd = "cp $< $@",
    )

    native.genrule(
        name = dname,
        srcs = [descriptor],
        outs = [name + ".descriptor_"],
        cmd = "cp $< $@",
    )

    _grpcstar_entrypoint(
        name = goname,
        main = starname,
        descriptor = dname,
    )

    go_library(
        name = libname,
        srcs = srcs + [goname],
        importpath = importpath,
        embedsrcs = [starname, dname] + scripts,
        visibility = visibility,
        deps = deps + [
            str(Label("//pkg/program")),
            str(Label("//pkg/protodescriptorset")),
        ],
    )

    go_binary(
        name = name,
        embed = [libname],
        visibility = visibility,
        **kwargs
    )
