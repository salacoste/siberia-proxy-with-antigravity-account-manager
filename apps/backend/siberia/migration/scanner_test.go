package migration

import (
	"testing"
)

// TestScanIDEProfiles is a functional test that runs the scanner on the current system
func TestScanIDEProfiles(t *testing.T) {
	profiles, err := ScanIDEProfiles()
	if err != nil {
		t.Fatalf("Scanner failed: %v", err)
	}

	t.Logf("Found %d profiles", len(profiles))
	for _, p := range profiles {
		t.Logf("- Found: %s (%s)", p.Name, p.Version)
		t.Logf("  Config: %s", p.ConfigPath)
		t.Logf("  DB: %s", p.DBPath)
		t.Logf("  Process: %s", p.ProcessName)

		// Optional: Try to read if DB exists
		// val, err := ReadStateToken(p.DBPath)
		// if err != nil {
		// 	t.Logf("  [WARN] Read failed: %v", err)
		// } else {
		// 	t.Logf("  [INFO] Token Length: %d", len(val))
		// }
	}

	// Not failing if 0 found, as this depends on the user's machine environment.
	// But it verifies the code doesn't crash.
}
