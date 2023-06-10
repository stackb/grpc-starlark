def main(ctx):
    print("=== Example Time Usage ===")
    now = time.now()
    then = now + 5 * time.second

    print("time then:", then)
    print("time hours:", then.hour)
    print("time minute:", then.minute)
    print("time second:", then.second)
    print("additional details: https://github.com/google/starlark-go/blob/a134d8f9ddca7469c736775b67544671f0a135ad/starlark/testdata/time.star")
