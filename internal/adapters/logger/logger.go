package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/sangrita-tech/periscope/internal/config"
)

const (
	timestampField = "timestamp"
	devTimeFormat  = "2006-01-02 15:04:05"
)

func New(cfg *config.Logger) (*zerolog.Logger, error) {
	if cfg == nil {
		cfg = &config.Logger{}
	}

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	zerolog.TimestampFieldName = timestampField
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	output := io.Writer(os.Stderr)
	if cfg.DevMode {
		output = &devWriter{output: output}
	}

	ctx := zerolog.New(output).Level(level).With().Timestamp()
	for _, key := range sortedKeys(cfg.BaseFields) {
		ctx = ctx.Str(key, cfg.BaseFields[key])
	}

	logger := ctx.Logger()
	return &logger, nil
}

type devWriter struct {
	output io.Writer
	mu     sync.Mutex
}

func (w *devWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for line := range bytes.SplitSeq(bytes.TrimSuffix(data, []byte("\n")), []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		formatted, err := formatDevLine(line)
		if err != nil {
			if _, writeErr := w.output.Write(append(line, '\n')); writeErr != nil {
				return 0, writeErr
			}
			continue
		}

		if _, err := io.WriteString(w.output, formatted+"\n"); err != nil {
			return 0, err
		}
	}

	return len(data), nil
}

func formatDevLine(line []byte) (string, error) {
	event := make(map[string]any)
	if err := json.Unmarshal(line, &event); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s %-5s | %s%s",
		formatTime(event[timestampField]),
		strings.ToUpper(asString(event[zerolog.LevelFieldName])),
		asString(event[zerolog.MessageFieldName]),
		formatFields(event),
	), nil
}

func formatTime(value any) string {
	timestamp, err := time.Parse(zerolog.TimeFieldFormat, asString(value))
	if err != nil {
		return asString(value)
	}
	return timestamp.UTC().Format(devTimeFormat)
}

func formatFields(event map[string]any) string {
	keys := make([]string, 0, len(event))
	for key := range event {
		switch key {
		case timestampField, zerolog.LevelFieldName, zerolog.MessageFieldName:
			continue
		default:
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, jsonString(key)+": "+jsonValue(event[key]))
	}

	return " {" + strings.Join(parts, ", ") + "}"
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func jsonString(value string) string {
	data, err := json.Marshal(value)
	if err != nil {
		return `""`
	}
	return string(data)
}

func jsonValue(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return `null`
	}
	return string(data)
}

func asString(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return fmt.Sprintf("%v", value)
}

func parseLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case "", "info":
		return zerolog.InfoLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "warn", "warning":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	default:
		return zerolog.NoLevel, fmt.Errorf("unsupported logger level %q", level)
	}
}
