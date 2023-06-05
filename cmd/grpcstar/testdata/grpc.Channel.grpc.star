# grpc.Channel is used by a grpc.Client
print("grpc.Channel:", grpc.Channel)
print("dir(grpc.Channel):", dir(grpc.Channel))
channel = grpc.Channel("localhost:9654")
print("channel:", channel)
