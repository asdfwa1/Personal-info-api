package logger

import (
	"Test_Task_EffMob/internal/middleware"
	"context"
	"log/slog"
	"os"
)

type contextHandler struct {
	slog.Handler
}

var (
	dfltLogger *slog.Logger
)

func InitLogger(level slog.Level) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	ctxHandler := &contextHandler{Handler: handler}
	dfltLogger = slog.New(ctxHandler)
	slog.SetDefault(dfltLogger)
}

func (ch *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID := middleware.GetRequestID(ctx); requestID != "" {
		r.AddAttrs(slog.String("requestID", requestID))
	}
	return ch.Handler.Handle(ctx, r)
}
