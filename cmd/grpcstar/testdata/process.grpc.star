def main(ctx):
    print(process)
    print("executable:", process.executable)
    print("run:", process.run)

    result = process.run(
        command = "awk",
        # args = ["--help", "all"],
    )
    print("stdout:", result.stdout)
    print("stderr:", result.stderr)
    print("error:", result.error)
    print("exit_code:", result.exit_code)
    print(result)
