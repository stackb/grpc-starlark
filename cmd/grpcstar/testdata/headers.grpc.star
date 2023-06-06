"""headers.grpc.star 

Demonstrates the sending and revcing of headers for unary and streaming methods.

"""
pb = proto.package("example.routeguide")
service_name = "example.routeguide.RouteGuide"

# === [Server Handler Functions] ================================================

def get_feature(stream, point):
    """get_feature implements a unary method handler

    Args:
        stream: the stream object
        point: the requested Point
    Returns:
        a Feature, ideally nearest to the given point.

    """
    md = stream.ctx.metadata()
    print("server: GetFeature request message:", point)
    print("server: GetFeature request headers:", dir(md))
    print("server: GetFeature request header content-type:", getattr(md, "content-type", "NOT SET"))
    print("server: GetFeature request header user-agent:", getattr(md, "user-agent", "NOT SET"))
    print("server: GetFeature request header x-unary-request:", getattr(md, "x-unary-request", "NOT SET"))

    stream.set_header({
        "x-unary-response": "pong!",
    })
    return pb.Feature(name = "point (%d,%d)" % (point.longitude, point.latitude))

def record_route(stream):
    """record_route implements a client streaming handler

    Args:
        stream: the stream object
    Returns:
        a RouteSummary with a summary of the traversed points.

    """
    md = stream.ctx.metadata()
    print("server: RecordRoute request headers:", dir(md))
    print("server: RecordRoute request header content-type:", getattr(md, "content-type", "NOT SET"))
    print("server: RecordRoute request header user-agent:", getattr(md, "user-agent", "NOT SET"))
    print("server: RecordRoute request header x-streaming-request:", getattr(md, "x-streaming-request", "NOT SET"))

    stream.set_header({
        "x-streaming-response": "pong!",
    })

    points = []
    for point in stream:
        points.append(point)

    stream.set_trailer({
        "x-streaming-trailer": "done!",
    })

    return pb.RouteSummary(
        point_count = len(points),
    )

listener = net.Listener()
server = grpc.Server()
server.register(service_name, {
    "GetFeature": get_feature,
    "RecordRoute": record_route,
})
thread.defer(lambda: server.start(listener))

channel = grpc.Channel(listener.address)
client = grpc.Client(service_name, channel)

def call_get_feature():
    feature = client.GetFeature(
        request = pb.Point(longitude = 1, latitude = 2),
        metadata = {
            "x-unary-request": "ping!",
        },
    )
    print("client: GetFeature response message:", feature)

def call_record_route():
    md = grpc.Metadata()
    md["x-streaming-request"] = "ping!"
    stream = client.RecordRoute(md)
    stream.close_send()

    route = stream.recv()
    print("client: RecordRoute response message:", route)

    headers = stream.header()
    print("client: RecordRoute response headers:", headers)
    print("client: RecordRoute response header x-streaming-response:", getattr(headers, "x-streaming-response", "NOT SET"))

    trailers = stream.trailer()
    print("client: RecordRoute response trailers:", dir(trailers))
    print("client: RecordRoute response trailer x-streaming-trailer:", getattr(trailers, "x-streaming-trailer", "NOT SET"))

def main(ctx):
    call_get_feature()
    call_record_route()

    server.stop()
