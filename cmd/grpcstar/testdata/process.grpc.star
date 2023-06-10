def main(ctx):
    print(process)

    print("run:", process.run)
    result = process.run(
        command = "pwd",
    )
    print("stdout (runfiles dir):", str(result.stdout).partition("grpcstar_test.runfiles")[2])
    print("stderr:", result.stderr)
    print("error:", result.error)
    print("exit_code:", result.exit_code)
