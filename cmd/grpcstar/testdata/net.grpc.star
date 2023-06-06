def main(ctx):
    print("=== Example Net Usage ===")
    print("net.Listener constructor", net.Listener)
    listener = net.Listener(network = "tcp", address = "localhost:1301")
    print("listener instance:", listener)
    print("listener.address:", listener.address)
    print("listener.network:", listener.network)
