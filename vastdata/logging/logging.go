// Package logging is the single logging entry point of the provider's resource logic.
package logging

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Logger is the logging contract hosts implement to receive the provider's logs.
// additionalFields are merged key/value pairs, exactly like tflog's.
type Logger interface {
	Trace(ctx context.Context, msg string, additionalFields ...map[string]any)
	Debug(ctx context.Context, msg string, additionalFields ...map[string]any)
	Info(ctx context.Context, msg string, additionalFields ...map[string]any)
	Warn(ctx context.Context, msg string, additionalFields ...map[string]any)
	Error(ctx context.Context, msg string, additionalFields ...map[string]any)
}

// TFLogger is the default Logger: it forwards to tflog.
type TFLogger struct{}

var _ Logger = TFLogger{}

func (TFLogger) Trace(ctx context.Context, msg string, f ...map[string]any) {
	tflog.Trace(ctx, msg, f...)
}
func (TFLogger) Debug(ctx context.Context, msg string, f ...map[string]any) {
	tflog.Debug(ctx, msg, f...)
}
func (TFLogger) Info(ctx context.Context, msg string, f ...map[string]any) {
	tflog.Info(ctx, msg, f...)
}
func (TFLogger) Warn(ctx context.Context, msg string, f ...map[string]any) {
	tflog.Warn(ctx, msg, f...)
}
func (TFLogger) Error(ctx context.Context, msg string, f ...map[string]any) {
	tflog.Error(ctx, msg, f...)
}

var (
	mu     sync.RWMutex
	global Logger = TFLogger{}
)

// SetLogger installs the process-wide logger; nil restores the tflog default.
// Call it once at start-up. A logger attached to a context takes precedence.
func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	if l == nil {
		l = TFLogger{}
	}
	global = l
}

type ctxKey struct{}

// WithLogger returns a context whose logging calls go to l.
func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the Logger that calls with ctx are routed to.
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		if l, ok := ctx.Value(ctxKey{}).(Logger); ok && l != nil {
			return l
		}
	}
	mu.RLock()
	defer mu.RUnlock()
	return global
}

func Trace(ctx context.Context, msg string, f ...map[string]any) {
	FromContext(ctx).Trace(ctx, msg, f...)
}
func Debug(ctx context.Context, msg string, f ...map[string]any) {
	FromContext(ctx).Debug(ctx, msg, f...)
}
func Info(ctx context.Context, msg string, f ...map[string]any) {
	FromContext(ctx).Info(ctx, msg, f...)
}
func Warn(ctx context.Context, msg string, f ...map[string]any) {
	FromContext(ctx).Warn(ctx, msg, f...)
}
func Error(ctx context.Context, msg string, f ...map[string]any) {
	FromContext(ctx).Error(ctx, msg, f...)
}
