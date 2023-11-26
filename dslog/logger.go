package dslog

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

var logger *slog.Logger

const (
	LevelDebug slog.Level = -4
	LevelInfo  slog.Level = 0
	LevelWarn  slog.Level = 4
	LevelError slog.Level = 8
	LevelFatal slog.Level = 12
	LevelAudit slog.Level = 16
)

// Initialize initializes the logger.  This should be called at the start of your
// program.  Once initialized, the logger can be used anywhere in your program.
func Configure(option Option) {
	handler := option.NewDSHandler()
	logger = slog.New(handler)

	hostname, _ := os.Hostname()
	env := slog.Group("env",
		slog.Int("pid", os.Getpid()),
		slog.String("hostname", hostname),
	)
	logger = logger.With(env)
}

func init() {
	if logger == nil {
		Configure(Option{})
		slog.SetDefault(logger)
	}
}

// Debug logs a message at debug level.
func Debug(msg string, args ...interface{}) {
	Log(nil, LevelDebug, msg, args...)
}

// Info logs a message at info level.
func Info(msg string, args ...interface{}) {
	Log(nil, LevelInfo, msg, args...)
}

// Warn logs a message at warn level.
func Warn(msg string, args ...interface{}) {
	Log(nil, LevelWarn, msg, args...)
}

// Error logs a message at error level.
func Error(msg string, args ...interface{}) {
	Log(nil, LevelError, msg, args...)
}

// Fatal logs a message at fatal level.
func Fatal(msg string, args ...interface{}) {
	Log(nil, LevelFatal, msg, args...)
}

// Audit logs a message at audit level.  To be used for auditable events.
func Audit(msg string, args ...interface{}) {
	Log(nil, LevelAudit, msg, args...)
}

// DebugCtx logs a message at debug level with context.
func DebugCtx(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, LevelDebug, msg, args...)
}

// InfoCtx logs a message at info level with context.
func InfoCtx(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, LevelInfo, msg, args...)
}

// WarnCtx logs a message at warn level with context.
func WarnCtx(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, LevelWarn, msg, args...)
}

// ErrorCtx logs a message at error level with context.
func ErrorCtx(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, LevelError, msg, args...)
}

// FatalCtx logs a message at fatal level with context.
func FatalCtx(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, LevelFatal, msg, args...)
}

// AuditCtx logs a message at audit level with context.  To be used for auditable events.
func AuditCtx(ctx context.Context, msg string, args ...interface{}) {
	Log(ctx, LevelAudit, msg, args...)
}

// Log logs a message at the given level with context.
func Log(ctx context.Context, level slog.Level, msg string, args ...interface{}) {
	logger.Log(ctx, level, msg, args...)
}
