package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Игнорируем запись журнала
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращаем тот же обработчик т.к. нет атрибутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Возвращаем тот же обработчик т.к. нет группы для сохранения
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Возвращаем false т.к. запись в журнал игнориурется
	return false
}
