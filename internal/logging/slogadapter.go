package logging

import (
	"context"
	"log/slog"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(logger *slog.Logger) *SlogAdapter {
	return &SlogAdapter{logger: logger}
}

func (l *SlogAdapter) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogAdapter) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SlogAdapter) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SlogAdapter) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SlogAdapter) With(args ...any) Logger {
	return &SlogAdapter{logger: l.logger.With(args...)}
}
