def main(ctx):
    print("=== Example Time Usage ===")
    now = time.now()

    # print("current time:", now)
    print("time add:", now + 5 * time.second)
    print("time hours:", now.hour)
    print("additional details: https://github.com/google/starlark-go/blob/a134d8f9ddca7469c736775b67544671f0a135ad/starlark/testdata/time.star")
