package processlog

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestReaderParsesFatalRuntimeLogs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth-service.log")
	content := "2026/05/24 20:10:30.123456 auth-service listening on :8080\n" +
		"2026/05/24 20:11:30.123456 fatal: database unavailable\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	entries, err := Reader{Path: path}.ListRuntimeLogs(t.Context(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected two entries, got %d", len(entries))
	}
	if entries[1].Level != "fatal" || entries[1].Message != "fatal: database unavailable" {
		t.Fatalf("unexpected fatal entry: %#v", entries[1])
	}
}

func TestConfigureTagsStandardLogFatal(t *testing.T) {
	if os.Getenv("PROCESSLOG_FATAL_CHILD") == "1" {
		closeLog, err := Configure(os.Getenv("PROCESSLOG_FATAL_PATH"))
		if err != nil {
			panic(err)
		}
		_ = closeLog
		log.Fatal(errors.New("database unavailable"))
	}

	path := filepath.Join(t.TempDir(), "auth-service.log")
	cmd := exec.Command(os.Args[0], "-test.run=TestConfigureTagsStandardLogFatal")
	cmd.Env = append(os.Environ(), "PROCESSLOG_FATAL_CHILD=1", "PROCESSLOG_FATAL_PATH="+path)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected child process to exit after log.Fatal")
	}

	entries, err := Reader{Path: path}.ListRuntimeLogs(t.Context(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one entry, got %d", len(entries))
	}
	if entries[0].Level != "fatal" || entries[0].Message != "fatal: database unavailable" {
		t.Fatalf("unexpected log.Fatal entry: %#v", entries[0])
	}
}
