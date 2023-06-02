package program

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	starlarkproto "go.starlark.net/lib/proto"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
)

//go:embed routeguide_proto_descriptor.pb
var routeguideProtoDescriptor []byte

//go:embed routeguide.grpc.star
var routeguideServer string

func TestProgram(t *testing.T) {

	pds, err := protodescriptorset.Parse(routeguideProtoDescriptor)
	if err != nil {
		t.Fatal(err)
	}
	files, err := protodesc.NewFiles(pds)
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]struct {
		program     string
		env         map[string]string
		wantErr     string
		wantPrinted string
	}{
		"simple case": {
			program:     "print(grpc.status.OK)",
			wantPrinted: "0",
		},
		"proto.package": {
			program: `
pb = proto.package("example.routeguide")
print(pb)
`,
			wantPrinted: `<proto.Package "example.routeguide">`,
		},
		"proto.MessageType": {
			program: `
pb = proto.package("example.routeguide")
print(pb.Point)
`,
			wantPrinted: `<proto.MessageType "example.routeguide.Point">`,
		},
		"skycfg-services-not-included": {
			program: `
pb = proto.package("example.routeguide")
print(pb.RouteGuide)
`,
			// this test demonstrates that services are not included in the proto
			// package from skycfg
			wantErr: `Protobuf type "example.routeguide.RouteGuide" not found`,
		},
		"unary method": {
			program: routeguideServer + `
def call_get_feature(client):
	point = pb.Point(longitude = 1, latitude = 2)
	got = client.GetFeature(point)
	print("GetFeature:", got)
	server.stop()
	
call_get_feature(client)
`,
			wantPrinted: `
GetFeature: <example.routeguide.Feature name:"point (1,2)">
`,
		},
		"server streaming": {
			program: routeguideServer + `
def call_list_features(client):
	rect = pb.Rectangle(
		lo = pb.Point(longitude = 1, latitude = 2),
		hi = pb.Point(longitude = 3, latitude = 4),
	)
	stream = client.ListFeatures(rect)
	for response in stream:
		print("ListFeatures:", response)
	server.stop()

call_list_features(client)
		`,
			wantPrinted: `
ListFeatures: <example.routeguide.Feature name:"lo (1,2)">
ListFeatures: <example.routeguide.Feature name:"hi (1,4)">
`,
		},
		"client streaming": {
			program: routeguideServer + `
def call_record_route(client):
	stream = client.RecordRoute()
	stream.send(pb.Point(longitude = 1, latitude = 2))
	stream.send(pb.Point(longitude = 3, latitude = 3))
	stream.close_send()
	response = stream.recv()
	print("RecordRoute:", response)
	server.stop()

call_record_route(client)
		`,
			wantPrinted: `
RecordRoute: <example.routeguide.RouteSummary point_count:2 distance:2 elapsed_time:10>
`,
		},
		"bidi streaming": {
			program: routeguideServer + `
def call_route_chat(client):
	stream = client.RouteChat()
	stream.send(pb.RouteNote(message = 'A'))
	stream.send(pb.RouteNote(message = 'B'))
	stream.close_send()
	for response in stream:
		print("RouteChat:", response)
	server.stop()

call_route_chat(client)
		`,
			wantPrinted: `
RouteChat: <example.routeguide.RouteNote message:"A">
RouteChat: <example.routeguide.RouteNote message:"B">
`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// address is set by the starlark program by printing a magic string
			// like '!address ADDRESS'
			var address string

			var printed bytes.Buffer
			thread := new(starlark.Thread)
			thread.Name = "main:" + name
			thread.Print = func(thread *starlark.Thread, msg string) {
				t.Log(msg)
				fmt.Println(msg)
				if strings.HasPrefix(msg, "!address") {
					fields := strings.Fields(address)
					address = fields[1]
					return
				}
				printed.WriteString(msg)
				printed.WriteString("\n")
			}
			// thread.Load = func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			// 	if module == "routeguide_server.grpc.star" {

			// 	}
			// }
			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			globals := newPredeclared(files)
			starlarkproto.SetPool(thread, files)

			_, err = starlark.ExecFile(
				thread,
				"<in-memory>",
				strings.NewReader(tc.program),
				globals,
			)

			if err != nil {
				if tc.wantErr == "" {
					t.Error("unexpected error: ", err)
					return
				}
				gotErr := err.Error()
				if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
					t.Errorf("error (-want +got):\n%s", diff)
				}
				return
			}

			gotPrinted := strings.TrimSpace(printed.String())
			wantPrinted := strings.TrimSpace(tc.wantPrinted)

			if diff := cmp.Diff(wantPrinted, gotPrinted); diff != "" {
				t.Errorf("print (-want +got):\n%s", diff)
			}
		})
	}
}
