pb = proto.package("google.protobuf.compiler")

fakeRequest = pb.CodeGeneratorRequest()

def generate(request):
    return [pb.CodeGeneratorResponse()]

def main(ctx):
    # return generate(ctx.vars.request)
    return generate(fakeRequest)
