package sbdb

import (
	"context"
	"log/slog"
)

type nopHandler struct{}

func (nopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (nopHandler) Handle(context.Context, slog.Record) error { return nil }
func (nopHandler) WithAttrs([]slog.Attr) slog.Handler        { return nopHandler{} }
func (nopHandler) WithGroup(string) slog.Handler             { return nopHandler{} }

var log = slog.New(nopHandler{})

// SetLogger replaces the package-level logger. Passing nil disables logging.
func SetLogger(l *slog.Logger) {
	if l == nil {
		log = slog.New(nopHandler{})
		return
	}
	log = l
}

func logFailedTypeAssert(fn, field Field, value any) {
	log.Debug("Type assertion failed", "fn", fn, "field", field, "value", value)
}
