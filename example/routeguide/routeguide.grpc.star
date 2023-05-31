"""mocks.star contains handler definitions for an example grpc service. 
"""
routeguidepb = proto.package("example.routeguide")

def decode_feature(json_str):
    return proto.decode_json(routeguidepb.Feature, json_str)

def decode_feature_database(json_str):
    return proto.decode_json(routeguidepb.FeatureDatabase, json_str)

def get_feature(_point, _context):
    """get_feature implements a unary method handler

    Args:
        _point: the requested Point
        _context: the method context object
    Returns:
        a Feature, ideally nearest to the given point.

    """

    # return decode_feature("""
    # {
    #     "location": {
    #         "latitude": 407838351,
    #         "longitude": -746143763
    #     },
    #     "name": "Patriots Path, Mendham, NJ 07945, USA"
    # }
    # """)
    return routeguidepb.Feature(
        location = routeguidepb.Point(
            latitude = 407838351,
            longitude = -746143763,
        ),
        name = "Patriots Path, Mendham, NJ 07945, USA",
    )

def list_features(_rectangle, context):
    """list_features implements a server streaming handler

    Args:
        _rectangle: the rectangle to get features within
        context: the stream context object
    Returns:
        None

    """
    db = decode_feature_database("""
{
    "feature": [
        {
            "location": {
                "latitude": 407838351,
                "longitude": -746143763
            },
            "name": "Patriots Path, Mendham, NJ 07945, USA"
        },
        {
            "location": {
                "latitude": 408122808,
                "longitude": -743999179
            },
            "name": "101 New Jersey 10, Whippany, NJ 07981, USA"
        },
        {
            "location": {
                "latitude": 413628156,
                "longitude": -749015468
            },
            "name": "U.S. 6, Shohola, PA 18458, USA"
        }
    ]
}
    """)
    for feature in db.feature:
        context.send(feature)

def record_route(_, context):
    """record_route implements a client streaming handler

    Args:
        _: the request object, which in this case is None
        context: the stream context object
    Returns:
        a RouteSummary with a summary of the traversed points.

    """
    points = []
    for point in context.recv:
        points.append(point)
    return routeguidepb.RouteSummary(
        point_count = len(points),
    )

def route_chat(_, context):
    """route_chat implements a bidirectional streaming handler

    Args:
        _: the request object, which in this case is None
        context: the stream context object
    Returns:
        None

    """
    notes = []
    for note in context.recv:
        notes.append(note)
        context.send(note)

grpc.RegisterHandlers({
    "/example.routeguide.RouteGuide/GetFeature": get_feature,
    "/example.routeguide.RouteGuide/ListFeatures": list_features,
    "/example.routeguide.RouteGuide/RecordRoute": record_route,
    "/example.routeguide.RouteGuide/RouteChat": route_chat,
})
