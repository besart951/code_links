package processlog

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/domain"
)

const DefaultPath = "logs/auth-service.log"

var (
	mu         sync.Mutex
	configured bool
)

func Configure(path string) (func() error, error) {
	mu.Lock()
	defer mu.Unlock()

	if configured {
		return func() error { return nil }, nil
	}
	if path == "" {
		path = DefaultPath
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
	log.SetOutput(fatalTaggingWriter{writers: []io.Writer{os.Stderr, file}})
	configured = true

	return file.Close, nil
}

func ConfigureFromEnv() (func() error, error) {
	return Configure(os.Getenv("AUTH_SERVICE_LOG_FILE"))
}

type Reader struct {
	Path string
}

func (r Reader) ListRuntimeLogs(_ context.Context, limit int) ([]domain.RuntimeLogEntry, error) {
	path := r.Path
	if path == "" {
		path = DefaultPath
	}
	if limit < 1 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return []domain.RuntimeLogEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines, err := tailLines(file, limit)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.RuntimeLogEntry, 0, len(lines))
	for _, line := range lines {
		if entry, ok := parseLine(line); ok {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func tailLines(file *os.File, limit int) ([]string, error) {
	const maxBytes int64 = 2 << 20

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	offset := int64(0)
	if stat.Size() > maxBytes {
		offset = stat.Size() - maxBytes
	}
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	if offset > 0 && scanner.Scan() {
		// Drop partial first line from the middle of the file.
	}

	lines := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if len(lines) > limit {
			lines = lines[len(lines)-limit:]
		}
	}
	return lines, scanner.Err()
}

func parseLine(line string) (domain.RuntimeLogEntry, bool) {
	const prefixLength = len("2006/01/02 15:04:05.000000")
	if len(line) <= prefixLength+1 {
		return domain.RuntimeLogEntry{}, false
	}

	occurredAt, err := time.Parse("2006/01/02 15:04:05.000000", line[:prefixLength])
	if err != nil {
		return domain.RuntimeLogEntry{}, false
	}
	message := strings.TrimSpace(line[prefixLength:])
	level := "info"
	if strings.HasPrefix(strings.ToLower(message), "fatal:") {
		level = "fatal"
	}

	hash := sha1.Sum([]byte(line))
	return domain.RuntimeLogEntry{
		ID:         hex.EncodeToString(hash[:]),
		OccurredAt: occurredAt,
		Level:      level,
		Source:     "auth-service",
		Message:    message,
		Raw:        line,
	}, true
}

type fatalTaggingWriter struct {
	writers []io.Writer
}

func (w fatalTaggingWriter) Write(p []byte) (int, error) {
	out := p
	if calledFromLogFatal() {
		out = tagFatal(p)
	}

	for _, writer := range w.writers {
		if _, err := writer.Write(out); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

func calledFromLogFatal() bool {
	var pcs [16]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		name := frame.Function
		if name == "log.Fatal" || name == "log.Fatalf" || name == "log.Fatalln" ||
			strings.HasSuffix(name, ".(*Logger).Fatal") ||
			strings.HasSuffix(name, ".(*Logger).Fatalf") ||
			strings.HasSuffix(name, ".(*Logger).Fatalln") {
			return true
		}
		if !more {
			return false
		}
	}
}

func tagFatal(p []byte) []byte {
	line := string(p)
	const prefixLength = len("2006/01/02 15:04:05.000000")
	if len(line) > prefixLength+1 && line[prefixLength] == ' ' {
		message := strings.TrimLeft(line[prefixLength+1:], " \t")
		if strings.HasPrefix(strings.ToLower(message), "fatal:") {
			return p
		}
		return []byte(line[:prefixLength+1] + "fatal: " + message)
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "fatal:") {
		return p
	}
	return []byte("fatal: " + line)
}
