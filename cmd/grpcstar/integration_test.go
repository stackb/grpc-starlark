package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	libtime "go.starlark.net/lib/time"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldens(t *testing.T) {
	os.Setenv("GODEBUG", "http2debug=2")
	flag.Parse()
	workspaceDir := os.Getenv("BUILD_WORKING_DIRECTORY")

	start := time.Now()
	epoch := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	libtime.NowFunc = func() time.Time {
		delta := time.Since(start).Round(100 * time.Millisecond)
		now := epoch.Add(delta)
		return now
	}

	type goldenTest struct {
		file        string
		goldenFile  string
		logFilename string
	}
	var tests []*goldenTest

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no test files found!")
	}

	for _, file := range entries {
		if strings.HasSuffix(file.Name(), ".grpc.star") {
			goldenName := file.Name() + ".out"
			logFilename := file.Name() + ".log"
			tests = append(tests, &goldenTest{
				file:        file.Name(),
				goldenFile:  goldenName,
				logFilename: logFilename,
			})
		}
	}

	for _, pair := range tests {
		t.Run(pair.file, func(t *testing.T) {
			if err := run(".", []string{
				"-log_file=" + pair.logFilename,
				"-protoset=../../example/routeguide/routeguide_proto_descriptor.pb",
				filepath.Join("testdata", pair.file),
			}); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(pair.logFilename)
			if err != nil {
				t.Fatal("reading log file:", err)
			}
			if *update {
				if workspaceDir == "" {
					t.Fatal("BUILD_WORKING_DIRECTORY not set!")
				}
				srcFilename := filepath.Join(workspaceDir, "cmd", "grpcstar", "testdata", pair.goldenFile)
				if err := os.WriteFile(srcFilename, got, os.ModePerm); err != nil {
					t.Fatal("writing golden file:", err)
				}
			} else {
				want, err := os.ReadFile(filepath.Join("testdata", pair.goldenFile))
				if err != nil {
					t.Fatal("reading golden file:", err)
				}
				if diff := cmp.Diff(string(want), string(got)); diff != "" {
					t.Errorf("expr (-want +got):\n%s", diff)
				}
			}
		})
	}
}
