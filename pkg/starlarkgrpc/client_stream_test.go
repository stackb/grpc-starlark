package starlarkgrpc

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stackb/grpc-starlark/mocks"
	"go.starlark.net/starlark"
)

func init() {

}

type clientStreamTest struct {
	expr         string
	want         string
	wantPrefix   string
	wantErr      string
	clientStream *mocks.ClientStream
	value        starlark.Value
}

func (tc *clientStreamTest) Run(t *testing.T) {
	files := RouteguideFiles(t)
	md, ok := methodDescriptorByName(files,
		"example.routeguide.RouteGuide",
		"GetFeature",
	)
	if !ok {
		t.Fatal("method not found")
	}

	thread := new(starlark.Thread)

	mockClientStream := mocks.NewClientStream(t)
	mockClientStream.
		On("Context").
		Once().
		Return(context.Background())

	globals := starlark.StringDict{
		"stream": newClientStream(mockClientStream, md),
	}

	value, err := starlark.Eval(
		thread,
		"<expr>",
		tc.expr,
		globals,
	)

	if err != nil {
		if tc.wantErr == "" {
			t.Fatal("unexpected error: ", err)
		}
		gotErr := err.Error()
		if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
			t.Fatalf("(-want +got):\n%s", diff)
		}
		return
	}

	got := value.String()
	var want string
	if tc.want != "" {
		want = strings.TrimSpace(tc.want)
	} else if tc.wantPrefix != "" {
		want = tc.wantPrefix
		got = got[0:len(want)]
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("expr (-want +got):\n%s", diff)
	}
}

func TestStarlarkGrpcClientStreamExpr(t *testing.T) {
	for _, tc := range []clientStreamTest{
		{expr: "stream.recv", want: `<built-in function grpc.ClientStream.recv>`},
		{expr: "stream.send", want: `<built-in function grpc.ClientStream.send>`},
		{expr: "stream.close_send", want: `<built-in function grpc.ClientStream.close_send>`},
		{expr: "stream.header", want: `<built-in function grpc.ClientStream.header>`},
		{expr: "stream.trailer", want: `<built-in function grpc.ClientStream.trailer>`},
		{expr: "stream.ctx", want: `Context(metadata = <built-in function Context.metadata>, with_timeout = <built-in function Context.with_timeout>, with_value = <built-in function Context.with_value>)`},
		{expr: "stream.descriptor", wantPrefix: `protoreflect.MethodDescriptor(`},
	} {
		t.Run(tc.expr, func(t *testing.T) {
			tc.Run(t)
		})
	}
}

// func TestStarlarkGrpcClientStreamRecv(t *testing.T) {
// 	files := RouteguideFiles(t)
// 	md, ok := methodDescriptorByName(files,
// 		"example.routeguide.RouteGuide",
// 		"GetFeature",
// 	)
// 	if !ok {
// 		t.Fatal("method not found")
// 	}

// 	mockClientStream := mocks.NewClientStream(t)
// 	mockClientStream.
// 		On("Context").
// 		Once().
// 		Return(context.Background())

// 	globals := starlark.StringDict{
// 		"stream": newClientStream(mockClientStream, md),
// 	}

// 	tc := moduletest.ExecFileTest{
// 		Source: `
// stream.recv()
// `,
// 	}
// }
