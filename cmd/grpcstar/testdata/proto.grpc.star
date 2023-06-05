print("proto:", proto)
print("proto.package:", proto.package)
print("proto.package('example.routeguide'):", proto.package("example.routeguide"))
# print("proto.package('i-dont-exist'):", proto.package("i-dont-exist"))

pb = proto.package("example.routeguide")
print("dir(pb):", dir(pb))
