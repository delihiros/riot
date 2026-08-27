package main

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"
)

func TestCLIExitsNonZeroForMissingRequiredFlag(t *testing.T) {
	command := exec.Command("go", "run", ".", "--api-key", "dummy", "valorant", "match", "recent-matches")
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatal("expected a non-zero exit status")
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected an exit status error, got %T", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit status = %d, want 1", exitErr.ExitCode())
	}
	if !bytes.Contains(output, []byte(`required flag(s) "queue" not set`)) {
		t.Fatalf("unexpected output: %s", output)
	}
}
