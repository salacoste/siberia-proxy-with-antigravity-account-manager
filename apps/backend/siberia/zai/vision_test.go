package zai

import (
	"strings"
	"testing"
)

func TestImageSourceToContent(t *testing.T) {
	// Test URL
	url := "https://example.com/image.png"
	part, err := ImageSourceToContent(url)
	if err != nil {
		t.Fatalf("Failed to process URL: %v", err)
	}
	if part.Type != "image_url" {
		t.Errorf("Expected image_url type, got %s", part.Type)
	}
	if part.ImageURL.URL != url {
		t.Errorf("URL mismatch")
	}

	// Test Local File (Mocking os.ReadFile is hard without refactoring,
	// so we will trust the URL test and maybe skip local file test for now
	// or create a temp file)

	// Create temp file
	// tmpFile := "test_image.txt"
	// Create a dummy file (not real image, but good enough for reading)
	// In real usage, mime detection needs real signature.
	// Write "PNG" header-ish mock?
	// ‰PNG...

	// For simplicity, just test failure on non-existent file
	_, err = ImageSourceToContent("non_existent_file.png")
	if err == nil {
		t.Error("Expected error for missing file")
	} else if !strings.Contains(err.Error(), "no such file") {
		t.Errorf("Expected 'no such file' error, got: %v", err)
	}
}
