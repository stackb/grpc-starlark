"""routeguide.grpc.star contains a simple server implementation

Used in testing.
"""
pb = proto.package("example.routeguide")

def get_feature(point):
    """get_feature implements a unary method handler

    Args:
        point: the requested Point
    Returns:
        a Feature, ideally nearest to the given point.

    """
    return pb.Feature(name = "point (%d,%d)" % (point.longitude, point.latitude))

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

server = grpc.Server()

server.register("example.routeguide.RouteGuide", {
    "GetFeature": get_feature,
    "ListFeatures": list_features,
    "RecordRoute": record_route,
    "RouteChat": route_chat,
})

listener = net.Listener()
channel = grpc.Channel(listener.address)
client = grpc.Client("example.routeguide.RouteGuide", channel)

thread.timeout(lambda: server.start(listener))
