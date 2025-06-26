package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// Helper to create files and directories for testing
func createTestFiles(t *testing.T, root string, files map[string]string) {
	for path, content := range files {
		fullPath := filepath.Join(root, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write file %s: %v", fullPath, err)
		}
	}
}

func TestListFiles_NonRecursive(t *testing.T) {
	tmpDir := t.TempDir()
	files := map[string]string{
		"file1.yaml":        "a: b",
		"file2.yml":         "c: d",
		"file3.json":        "{}",
		"file4.txt":         "not yaml",
		"subdir/file5.yaml": "e: f",
	}
	createTestFiles(t, tmpDir, files)

	got, err := listFiles(tmpDir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		filepath.Join(tmpDir, "file1.yaml"),
		filepath.Join(tmpDir, "file2.yml"),
		filepath.Join(tmpDir, "file3.json"),
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestListFiles_Recursive(t *testing.T) {
	tmpDir := t.TempDir()
	files := map[string]string{
		"file1.yaml":        "a: b",
		"file2.yml":         "c: d",
		"file3.json":        "{}",
		"file4.txt":         "not yaml",
		"subdir/file5.yaml": "e: f",
		"subdir/file6.json": "{}",
		"subdir2/file7.yml": "g: h",
		"subdir2/file8.txt": "not yaml",
	}
	createTestFiles(t, tmpDir, files)

	got, err := listFiles(tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		filepath.Join(tmpDir, "file1.yaml"),
		filepath.Join(tmpDir, "file2.yml"),
		filepath.Join(tmpDir, "file3.json"),
		filepath.Join(tmpDir, "subdir", "file5.yaml"),
		filepath.Join(tmpDir, "subdir", "file6.json"),
		filepath.Join(tmpDir, "subdir2", "file7.yml"),
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestListFiles_Error(t *testing.T) {
	// Non-existent directory should return an error
	_, err := listFiles("/non/existent/dir", false)
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}
