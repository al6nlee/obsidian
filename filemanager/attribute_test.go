package filemanager_test

import (
	"github.com/al6nlee/obsidian/filemanager"
	"os"
	"testing"
	"time"
)

func TestProcessFiles(t *testing.T) {
	dir := "/Users/alan/Github/obsidian/note"
	err := filemanager.ProcessFiles(dir)
	if err != nil {
		t.Errorf("ProcessFiles() error = %v", err)
	}

	// You can add more assertions here if needed
}

func TestAddAttribute(t *testing.T) {
	file := filemanager.FileAttribute{
		Tag:        [2]string{"tag1", "tag2"},
		FileName:   "test.md",
		CreateTime: time.Now(),
		ModTime:    time.Now(),
		Size:       100,
		Mode:       "-rw-r--r--",
		Author:     "alan",
	}

	testFile := "/Users/alan/Github/obsidian/note/test.md"

	// Create a test file
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("error creating test file: %v", err)
	}
	defer f.Close()

	// Call AddAttribute
	err = filemanager.AddAttribute(testFile, file)
	if err != nil {
		t.Errorf("AddAttribute() error = %v", err)
	}

	// You can add more assertions here if needed
}
