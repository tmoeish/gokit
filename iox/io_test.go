package iox_test

import (
	"path/filepath"
	"testing"

	"github.com/tmoeish/gokit/iox"
)

func TestWriteReadFile(t *testing.T) {
	dir, cleanup, err := iox.TempDir("iox-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	path := filepath.Join(dir, "test.txt")
	content := []byte("hello gokit")

	if err := iox.WriteFile(path, content); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := iox.ReadFile(path)
	if err != nil || string(got) != string(content) {
		t.Fatalf("ReadFile: %s %v", got, err)
	}
}

func TestWriteReadLines(t *testing.T) {
	dir, cleanup, err := iox.TempDir("iox-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	path := filepath.Join(dir, "lines.txt")
	lines := []string{"line1", "line2", "line3"}

	if err := iox.WriteLines(path, lines); err != nil {
		t.Fatalf("WriteLines: %v", err)
	}

	got, err := iox.ReadLines(path)
	if err != nil || len(got) != 3 || got[0] != "line1" {
		t.Fatalf("ReadLines: %v %v", got, err)
	}
}

func TestFileExists(t *testing.T) {
	path, cleanup, err := iox.TempFile("iox-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	if !iox.FileExists(path) {
		t.Fatal("FileExists should be true")
	}
	if iox.FileExists(path + ".nonexistent") {
		t.Fatal("FileExists should be false for missing file")
	}
}

func TestCopyFile(t *testing.T) {
	dir, cleanup, err := iox.TempDir("iox-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "subdir", "dst.txt")

	if err := iox.WriteFile(src, []byte("copy me")); err != nil {
		t.Fatal(err)
	}
	if err := iox.CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}
	got, err := iox.ReadFileStr(dst)
	if err != nil || got != "copy me" {
		t.Fatalf("CopyFile content: %q %v", got, err)
	}
}

func TestLineCount(t *testing.T) {
	dir, cleanup, err := iox.TempDir("iox-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	path := filepath.Join(dir, "count.txt")
	iox.WriteLines(path, []string{"a", "b", "c", "d", "e"}) //nolint:errcheck

	n, err := iox.LineCount(path)
	if err != nil || n != 5 {
		t.Fatalf("LineCount: %d %v", n, err)
	}
}
