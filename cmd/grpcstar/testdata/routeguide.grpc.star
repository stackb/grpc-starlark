"""routeguide.grpc.star contains a simple server implementation

Used in testing.

The client test functions have scope access to the client and server.
Typical usage is to invoke the script, start the server in a separate thread via
thread.defer, wait for client return values, and stop the server upon success.

"""
pb = proto.package("example.routeguide")

# === [Server Handler Functions] ================================================

def get_feature(_stream, point):
    """get_feature implements a unary method handler

    Args:
        _stream: the stream object
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

listener = net.Listener()
server = grpc.Server()
server.register("example.routeguide.RouteGuide", {
    "GetFeature": get_feature,
    "ListFeatures": list_features,
    "RecordRoute": record_route,
    "RouteChat": route_chat,
})

channel = grpc.Channel(listener.address)
client = grpc.Client("example.routeguide.RouteGuide", channel)

# === [Client Call Functions] ================================================

def call_get_feature():
    point = pb.Point(longitude = 1, latitude = 2)
    feature = client.GetFeature(point)
    print("GetFeature:", feature)

def call_list_features():
    rect = pb.Rectangle(
        lo = pb.Point(longitude = 1, latitude = 2),
        hi = pb.Point(longitude = 3, latitude = 4),
    )
    stream = client.ListFeatures(rect)
    for response in stream:
        print("ListFeatures:", response)

def call_record_route():
    stream = client.RecordRoute()
    stream.send(pb.Point(longitude = 1, latitude = 2))
    stream.send(pb.Point(longitude = 3, latitude = 3))
    stream.close_send()
    response = stream.recv()
    print("RecordRoute:", response)

def call_route_chat():
    stream = client.RouteChat()
    stream.send(pb.RouteNote(message = "A"))
    stream.send(pb.RouteNote(message = "B"))
    stream.close_send()
    for response in stream:
        print("RouteChat:", response)

def main(ctx):
    thread.defer(lambda: server.start(listener))
    call_get_feature()
    call_list_features()
    call_record_route()
    call_route_chat()
    server.stop()
