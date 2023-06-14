[![CI](https://github.com/stackb/grpc-starlark/actions/workflows/ci.yaml/badge.svg)](https://github.com/stackb/grpc-starlark/actions/workflows/ci.yaml)

# grpc-starlark

<table border="0">
  <tr>
    <td><img src="https://user-images.githubusercontent.com/50580/141900696-bfb2d42d-5d2c-46f8-bd9f-06515969f6a2.png" height="120"/></td>
    <td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="120"/></td>
    <td><img src="https://static.vecteezy.com/system/resources/previews/007/038/145/non_2x/nightingale-singing-tune-song-bird-musical-notes-music-concept-icon-in-circle-round-black-color-illustration-flat-style-image-vector.jpg" height="120"/></td>
    <!-- image credit: https://www.vecteezy.com/vector-art/7038145-nightingale-singing-tune-song-bird-musical-notes-music-concept-icon-in-circle-round-black-color-vector-illustration-flat-style-image -->
  </tr>
  <tr>
    <td>protobuf</td>
    <td>grpc</td>
    <td>starlark</td>
  </tr>
</table>

`grpc-starlark` is a:
- library for embedding a gRPC-capable starlark interpreter,
- standalone binary `grpcstar` that executes starlark scripts.

> The author pronounces this as `grip-ster` (like "napster", but you can say it however you like).

`grpcstar` use cases include:

- replacement for `grpcurl` when calling gRPC services from the command line
- stand-in for `postman`
- testing gRPC backends
- mocking gRPC backends in integration tests

## Installation

Download a binary from the [releases
page](https://github.com/stackb/grpc-starlark/releases), or install from source:

```sh
go install github.com/stackb/grpc-starlark/cmd/grpcstar@latest
```

## Usage

```
usage: grpcstar [OPTIONS...] [ARGS...]

github:
	https://github.com/stackb/grpc-starlark

options:
	-h, --help [optional, false]
		show this help screen
	-p, --protoset [required]
		filename of proto descriptor set
	-f, --file [required]
		filename of entrypoint starlark script
		(conventionally named *.grpc.star)
	-e, --entrypoint [optional, "main"]
		name of function in global scope to invoke upon script start
	-o, --output [optional, "json", oneof "json|proto|text|yaml"]
		formatter for output protobufs returned by entrypoint function
	-i, --interactive [optional, false]
		start a REPL session (rather then exec the entrypoint)

example:
	$ grpcstar \
		-p routeguide.pb \
		-f routeguide.grpc.star \
		-e call_get_feature \
		longitude=35.0 latitude=109.1
```

## Bazel Usage

See [bazel rule documentation](rules/)

### Proto Descriptor Set

grpcstar requires a precompiled proto descriptor set via the `--protoset` (`-p`)
flag. This file defines the universe of message, enum, and service types that
can be used in your script.

This file can be generated by the protoc `--descriptor_set_out` flag and is used
by other tools in the protobuf/gRPC ecosystem (see
[grpcurl](https://github.com/fullstorydev/grpcurl#protoset-files)).

For bazel users, the
[proto_descriptor_set](https://github.com/bazelbuild/rules_proto/blob/master/proto/private/rules/proto_descriptor_set.bzl) rule
can be used to generate this file.

> Have multiple descriptor files?  You can simply `cat` multiple
> `FileDescriptorProto` binary files together!

### Script File

The script file `--file` (`-f`) is the entrypoint file executed by the embedded starlark interpreter. 

Use load statements (e.g. `load("filename{.star}", "symbol")`) to populate
additional symbols into the entrypoint file.

### Script Entrypoint

The script **must** contain a function named `main` that takes a single
positional argument `ctx` (e.g.`def main(ctx):`).  The `--entrypoint` (`-e`)
flag can be used to override this.

The `ctx` is a struct; `ctx.vars`  holds key-value pairs that can be set on the
command line (e.g. `name=foo` would satisfy `ctx.vars.name == 'foo'`).

### Script Output

The entrypoint function can either return nothing (`None`) or a list of protobuf
messages.  The messages will be printed to stdout and formatted according the
the `--output` flag (`-o`).  Choose one of `json`, `proto`, `text`, or `yaml`;
default is `json`.

`print(...)` statements are sent to stderr.

### Script Concurrency Model

The starlark interpreter starts a single `main` thread for the top-level
entrypoint file.  Each invocation of a `grpc.Server` handler callback function
is run concurrently in a new thread.  `thread.defer` callbacks also occur in a
new thread.

## API

`grpc-starlark` is implemented using go and has an API similar to `grpc-go`.

### Protobuf

The message and enum types are available via the `proto.package` function:

```py
pb = proto.package("example.routeguide")
print(pb.Rectangle)
```

These define "strongly-typed" structs for use in creating and interacting with
protobuf messages:

```py
colorado = pb.Rectangle(
    lo = pb.Point(latitude = 36.999, longitude = -109.045),
    hi = pb.Point(latitude = 40.979, longitude = -102.051),
)
```

For more details see
[github.com/stripe/skycfg](https://github.com/stripe/skycfg), which provides the
core protobuf functionality.

## gRPC

### Server

Use the `grpc.Server` constructor to make a new server.  Use the register
function to provide function implementations for the service handlers.  Example:

```py
server = grpc.Server()

server.register("example.routeguide.RouteGuide", {
    "GetFeature": get_feature,
    "ListFeatures": list_features,
    "RecordRoute": record_route,
    "RouteChat": route_chat,
})
```

Use a `net.Listener` to bind the server to a network address:

```py
listener = net.Listener(address = "localhost:8080")
```

To bind to a free port, use the defaults (`localhost` is the `host` and `0` is
the port)

```py
listener = net.Listener()
print(listener.address) # localhost:50234
```

#### Unary RPC

```py
def get_feature(stream, point):
    """get_feature implements a unary method handler

    Args:
        stream: the stream object
        point: the requested Point
    Returns:
        a Feature, ideally nearest to the given point.

    """
    return pb.Feature(name = "point (%d,%d)" % (point.longitude, point.latitude))
```

The `stream` object can be used to access incoming headers `stream.ctx.metadata`
or set outgoing headers/trailers (`stream.set_header`, `stream.set_trailer`).

The second positional argument is the request message.

The function should return an appropriate response message or a `grpc.Error`
using an status code and message (e.g. `return grpc.Error(code =
grpc.status.UNAUTHENTICATED, message = "authorization header is required"))`)

#### Server Streaming RPC

```py

def list_features(stream, rect):
    """list_features implements a server streaming handler

    Args:
        stream: the stream object
        rect: the rectangle to get features within
    Returns:
        None

    """
    features = [
        pb.Feature(name = "lo (%d,%d)" % (rect.lo.longitude, rect.lo.latitude)),
        pb.Feature(name = "hi (%d,%d)" % (rect.lo.longitude, rect.hi.latitude)),
    ]
    for feature in features:
        stream.send(feature)
```

The `stream.send` function is used to post response messages.

#### Client Streaming RPC

```py

def record_route(stream):
    """record_route implements a client streaming handler

    Args:
        stream: the stream object
    Returns:
        a RouteSummary with a summary of the traversed points.

    """
    points = []
    for point in stream:
        points.append(point)
    return pb.RouteSummary(
        point_count = len(points),
        distance = 2,
        elapsed_time = 10,
    )
```

The `stream` is an iterable that will call `.RevcMsg` until the stream has been closed by the client.  

Alternatively, the function `stream.recv` can be used to get a single message, or `None` if there are no more messages.

The return value of the function should return an appropriately typed message.

#### Bidirectional Streaming RPC

```py
def route_chat(stream):
    """route_chat implements a bidirectional streaming handler

    Args:
        stream: the stream object
    Returns:
        None

    """
    notes = []
    for note in stream:
        notes.append(note)
        stream.send(note)
```

In this implementation the function broadcasts a reponse on every request.

## time

The `time` module contains time-related functions.  For details, see <https://github.com/google/starlark-go/blob/master/lib/time/time.go>.

## os

The `os` module contains functions for interacting with the operating system.  

- `os.getenv("NAME")` returns the value of the environment variable `NAME` or `None` if not set.

See <https://github.com/stackb/grpc-starlark/tree/master/cmd/grpcstar/testdata> for details.

## thread

The `thread` module can be used to interact with the interpreter threading model.

- `thread.sleep(duration)` pauses the current thread.
- `thread.defer(fn, delay, count)` runs another function in a new thread after
  the given delay.  An optional `count` argument will repeat the callback
  invocation.  This function is akin to the javascript functions `setTimout` and
  `setInterval`.
- `thread.name` returns the name of the current thread.

Example:

```py
thread.defer(lambda: server.start(listener))`
```

## net

The `net` module contains network-related functions.

- `net.Listener` constructs a new listener via the
  [net.Listen](https://pkg.go.dev/net#Listen) func.


## process

The `process` module contains subprocess-related functions.

- `process.run` runs a subprocess.

