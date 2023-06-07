def main(ctx):
    print("proto:", proto)
    print("proto.package:", proto.package)
    print("proto.package('example.routeguide'):", proto.package("example.routeguide"))

    pb = proto.package("example.routeguide")
    print("dir(pb):", dir(pb))
