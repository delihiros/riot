package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfigReturnsErrorsAndLoadsExplicitFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "riot.yaml")
	if err := os.WriteFile(path, []byte("riot:\n  api_key: from-file\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := ReadConfig(path); err != nil {
		t.Fatal(err)
	}
	if got := GetConfig().GetString("riot.api_key"); got != "from-file" {
		t.Fatalf("api key = %q", got)
	}
}

func TestReadConfigAllowsMissingDefaultFile(t *testing.T) {
	oldHome, hadHome := os.LookupEnv("HOME")
	if err := os.Setenv("HOME", t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if hadHome {
			_ = os.Setenv("HOME", oldHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
	}()

	if err := ReadConfig(""); err != nil {
		t.Fatalf("missing default config should be optional: %v", err)
	}
}

func TestReadConfigLoadsAPIKeyFromEnvironment(t *testing.T) {
	oldHome, hadHome := os.LookupEnv("HOME")
	oldKey, hadKey := os.LookupEnv("RIOT_API_KEY")
	_ = os.Setenv("HOME", t.TempDir())
	_ = os.Setenv("RIOT_API_KEY", "from-environment")
	defer func() {
		if hadHome {
			_ = os.Setenv("HOME", oldHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
		if hadKey {
			_ = os.Setenv("RIOT_API_KEY", oldKey)
		} else {
			_ = os.Unsetenv("RIOT_API_KEY")
		}
	}()

	if err := ReadConfig(""); err != nil {
		t.Fatal(err)
	}
	if got := GetConfig().GetString("riot.api_key"); got != "from-environment" {
		t.Fatalf("api key = %q", got)
	}
}
