package dslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"golang.org/x/exp/slog"
)

func TestLog(t *testing.T) {
	// Create a buffer to capture the logs
	var b bytes.Buffer

	// Configure the logger options
	option := Option{
		Writer: &b,
		Level:  LevelDebug,
	}
	Configure(option)

	// Log a message with some additional fields
	Log(
		context.Background(),
		LevelDebug,
		"hello world",
		slog.String("name", "deepcode-ai"),
		slog.Int("age", 1),
	)

	// Unmarshal the logged message into a map for easy testing
	var got map[string]interface{}
	if err := json.Unmarshal(b.Bytes(), &got); err != nil {
		t.Error(err)
	}

	// Check that the log message has the expected fields and values
	if got["level"] != LevelDebug.String() {
		t.Error("expected debug level")
	}
	if got["message"] != "hello world" {
		t.Error("expected hello world message")
	}

	if ctxVal, ok := got["context"].(map[string]interface{}); ok {
		if ctxVal["name"] != "deepcode-ai" {
			t.Error("expected name to be deepcode-ai")
		}
		if fmt.Sprintf("%v", ctxVal["age"]) != "1" {
			t.Error("expected age to be 1")
		}
	} else {
		t.Error("expected context")
	}

	if env, ok := got["env"].(map[string]interface{}); ok {
		if env["pid"] == nil {
			t.Error("expected pid")
		}
		if env["hostname"] == nil {
			t.Error("expected hostname")
		}
	} else {
		t.Error("expected env")
	}
}
