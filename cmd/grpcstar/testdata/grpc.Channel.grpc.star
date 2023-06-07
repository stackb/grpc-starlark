def main(ctx):
    print("grpc.Channel:", grpc.Channel)
    channel = grpc.Channel("localhost:9654")
    print("channel:", channel)
