package starlarkgrpc

import (
	"testing"
)

func TestStarlarkGrpcClientStreamExpr(t *testing.T) {
	for name, tc := range map[string]struct {
		wantErr string
		want    string
	}{
		"degenerate": {},
	} {
		t.Run(name, func(t *testing.T) {
			t.Log(tc.want)
		})
	}
}
