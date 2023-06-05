package moduletest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// FileSpec specifies the content of a test file.
type FileSpec struct {
	// Path is a slash-separated path relative to the test directory. If Path
	// ends with a slash, it indicates a directory should be created
	// instead of a file.
	Path string

	// Symlink is a slash-separated path relative to the test directory. If set,
	// it indicates a symbolic link should be created with this path instead of a
	// file.
	Symlink string

	// Content is the content of the test file.
	Content string

	// NotExist asserts that no file at this path exists.
	// It is only valid in CheckFiles.
	NotExist bool
}

// CreateFiles creates a directory of test files. This is a more compact
// alternative to testdata directories. CreateFiles returns a canonical path
// to the directory and a function to call to clean up the directory
// after the test.
func CreateFiles(t *testing.T, files []FileSpec) (dir string, cleanup func()) {
	t.Helper()
	dir, err := ioutil.TempDir(os.Getenv("TEST_TEMPDIR"), "gazelle_test")
	if err != nil {
		t.Fatal(err)
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.NotExist {
			t.Fatalf("CreateFiles: NotExist may not be set: %s", f.Path)
		}
		path := filepath.Join(dir, filepath.FromSlash(f.Path))
		if strings.HasSuffix(f.Path, "/") {
			if err := os.MkdirAll(path, 0o700); err != nil {
				os.RemoveAll(dir)
				t.Fatal(err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			os.RemoveAll(dir)
			t.Fatal(err)
		}
		if f.Symlink != "" {
			if err := os.Symlink(f.Symlink, path); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := ioutil.WriteFile(path, []byte(f.Content), 0o600); err != nil {
			os.RemoveAll(dir)
			t.Fatal(err)
		}
	}

	return dir, func() { os.RemoveAll(dir) }
}

// CheckFiles checks that files in "dir" exist and have the content specified
// in "files". Files not listed in "files" are not tested, so extra files
// are allowed.
func CheckFiles(t *testing.T, dir string, files []FileSpec) {
	t.Helper()
	for _, f := range files {
		path := filepath.Join(dir, f.Path)

		st, err := os.Stat(path)
		if f.NotExist {
			if err == nil {
				t.Errorf("asserted to not exist, but does: %s", f.Path)
			} else if !os.IsNotExist(err) {
				t.Errorf("could not stat %s: %v", f.Path, err)
			}
			continue
		}

		if strings.HasSuffix(f.Path, "/") {
			if err != nil {
				t.Errorf("could not stat %s: %v", f.Path, err)
			} else if !st.IsDir() {
				t.Errorf("not a directory: %s", f.Path)
			}
		} else {
			want := strings.TrimSpace(f.Content)
			gotBytes, err := ioutil.ReadFile(filepath.Join(dir, f.Path))
			if err != nil {
				t.Errorf("could not read %s: %v", f.Path, err)
				continue
			}
			got := strings.TrimSpace(string(gotBytes))
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("%s diff (-want,+got):\n%s", f.Path, diff)
			}
		}
	}
}
