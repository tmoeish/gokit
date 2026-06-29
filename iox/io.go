// Package iox provides I/O utility functions.
package iox

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// ReadAll reads all bytes from r and returns them.
// Returns an error if reading fails.
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// ReadFile reads and returns the entire content of a file.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileStr reads and returns the entire content of a file as a string.
func ReadFileStr(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// WriteFile writes data to a file, creating it if it doesn't exist.
// The file is truncated before writing.
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

// WriteFileStr writes a string to a file.
func WriteFileStr(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// AppendFile appends data to a file, creating it if it doesn't exist.
func AppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	_, err = f.Write(data)
	return err
}

// ReadLines reads a file and returns its lines (without trailing newlines).
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteLines writes lines to a file, one per line.
func WriteLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	w := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := w.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

// FileExists reports whether the path exists and is a regular file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// DirExists reports whether the path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// PathExists reports whether the path exists (file or directory).
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureDir creates all directories in path if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// CopyFile copies the contents of src to dst, creating dst if necessary.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close() //nolint:errcheck

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close() //nolint:errcheck

	_, err = io.Copy(out, in)
	return err
}

// TempFile creates a temporary file with the given prefix and returns its path and a cleanup function.
func TempFile(prefix string) (string, func(), error) {
	f, err := os.CreateTemp("", prefix)
	if err != nil {
		return "", nil, err
	}
	path := f.Name()
	_ = f.Close()
	return path, func() { _ = os.Remove(path) }, nil
}

// TempDir creates a temporary directory and returns its path and a cleanup function.
func TempDir(prefix string) (string, func(), error) {
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", nil, err
	}
	return dir, func() { _ = os.RemoveAll(dir) }, nil
}

// LineCount returns the number of lines in a file.
func LineCount(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close() //nolint:errcheck
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

// FileSize returns the size of a file in bytes.
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
