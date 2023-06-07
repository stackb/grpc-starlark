def main(ctx):
    print("=== Example Module Usage ===")
    console = module(
        "console",
        log = lambda msg: print(msg),
    )
    print(console)

    console.log("console.log!")
