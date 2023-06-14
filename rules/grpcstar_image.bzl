load("@io_bazel_rules_docker//container:image.bzl", "container_image")
load("@io_bazel_rules_docker//container:layer.bzl", "container_layer")

def grpcstar_image(**kwargs):
    """grpcstar_image is a macro that produces a container image for grpcstar.

    Args:
        **kwargs: the keyword args dict
    Returns:
        None
    """
    name = kwargs.pop("name")
    starname = name + "._.star"
    descriptorname = name + ".d.layer"
    binname = name + "bin"
    starlayer = name + ".star.layer"
    descriptorlayer = name + ".descriptor.layer"
    binlayer = name + ".bin.layer"

    executable = kwargs.pop("executable", str(Label("//cmd/grpcstar")))
    base = kwargs.pop("base", "@go_image_base//image")
    layers = kwargs.pop("layers", [])
    layers += [binlayer, descriptorlayer, starlayer]

    descriptor = kwargs.pop("descriptor", "")
    if not descriptor:
        fail("grpcstar_image.descriptor is required")
    main = kwargs.pop("main", "")
    if not main:
        fail("grpcstar_image.main is required")

    native.genrule(
        name = name + "b",
        srcs = [executable],
        outs = [binname],
        cmd = "cp $< $@",
        executable = True,
    )
    native.genrule(
        name = name + "d",
        srcs = [descriptor],
        outs = [descriptorname],
        cmd = "cp $< $@",
    )
    native.genrule(
        name = name + "m",
        srcs = [main],
        outs = [starname],
        cmd = "cp $< $@",
    )

    container_layer(name = binlayer, files = [binname])
    container_layer(name = descriptorlayer, files = [descriptorname])
    container_layer(name = starlayer, files = [starname])

    container_image(
        name = name,
        base = base,
        layers = layers,
        cmd = None,
        launcher = None,
        entrypoint = [
            "/" + binname,
            "--protoset",
            descriptorname,
            "--file",
            starname,
        ],
        **kwargs
    )

    # container_push(
    #     name = "push",
    #     format = "Docker",
    #     image = ":image",
    #     registry = "us.gcr.io/bzlio-260121",
    #     repository = "bezel/robinhood",
    #     tag = "latest",
    # )
