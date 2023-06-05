"""headers.grpc.star demonstrates the sending and revcing of headers
"""
pb = proto.package("example.routeguide")
service_name = "example.routeguide.RouteGuide"

# === [Server Handler Functions] ================================================

def get_feature(stream, point):
    """get_feature implements a unary method handler

    Args:
        _stream: the stream object
        point: the requested Point
    Returns:
        a Feature, ideally nearest to the given point.

    """
    md = stream.ctx.metadata()
    print("keys:", dir(md))
    print("content-type:", getattr(md, "content-type", "NOT SET"))
    print("user-agent:", getattr(md, "user-agent", "NOT SET"))
    print("authorization:", getattr(md, "authorization", "NOT SET"))
    return pb.Feature(name = "point (%d,%d)" % (point.longitude, point.latitude))

listener = net.Listener()
server = grpc.Server()
server.register(service_name, {
    "GetFeature": get_feature,
})

channel = grpc.Channel(listener.address)
client = grpc.Client(service_name, channel)
thread.defer(lambda: server.start(listener))

md = grpc.Metadata()
md["authorization"] = "bearer foo"
client.GetFeature(
    request = pb.Point(longitude = 1, latitude = 2),
    metadata = md,
)
server.stop()
