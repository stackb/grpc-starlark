def main(ctx):
    now = time.now()

    print("""
    > thread.name prints the current thread name
    """)
    print(thread.name())

    print("""
    > thread.sleep pauses the current thread for the given duration
    """)
    print("thread.sleep:", thread.sleep)
    thread.sleep(duration = 200 * time.millisecond)
    print("sleep before:", now)
    print("sleep after:", time.now())

    print("""
    > thread.defer runs a function in a separate thread after a given delay
    """)
    thread.defer(
        fn = lambda: print("defer at %s in thread %s:" % (time.now(), thread.name())),
        delay = 100 * time.millisecond,
        count = 3,
    )

    # sleeping here only to get the print statements to line up better
    thread.sleep(500 * time.millisecond)

    print("""
    > thread.cancel stops the current thread.  Cancelling the main thread exits the program.
    """)

    print(thread.defer(fn = lambda: thread.cancel(reason = "example")))
