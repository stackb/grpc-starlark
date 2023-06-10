def main(ctx):
    print(process)

    print("run:", process.run)
    result = process.run(
        command = "ls",
    )
    print("stdout:", result.stdout)
    print("stderr:", result.stderr)
    print("error:", result.error)
    print("exit_code:", result.exit_code)
