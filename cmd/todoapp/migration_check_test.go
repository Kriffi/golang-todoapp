package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLatestMigrationVersionUsesNumericUpVersions(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"9_a.up.sql", "10_b.up.sql", "99_c.down.sql"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	version, err := latestMigrationVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if version != 10 {
		t.Fatalf("version = %d, want 10", version)
	}
}
