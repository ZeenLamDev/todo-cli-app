package main

import (
	"fmt"
	"log/slog"
	"os"

	"todo/app"
	"todo/logutil"
	"todo/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	ctx, cancel := app.SignalNotifyContext()
	traceID := fmt.Sprintf("trace-%d", os.Getpid())
	ctx = logutil.WithTraceID(ctx, traceID)

	app := &app.App{
		Todos:   store.NewTodos(),
		Storage: store.NewStorage[store.Todos]("todos.json"),
		Ctx:     ctx,
		Cancel:  cancel,
	}

	go app.StartHTTPServer("8080")

	if err := app.Storage.Load(app.Ctx, &app.Todos); err != nil {
		slog.Warn("Could not load todos. Starting with empty list", slog.Any("error", err))
	}

	slog.Info("🚀 App is running... Ctrl+C to stop")

	<-app.Ctx.Done()
	slog.Info("🛑 Interrupt received. Shutting down...")

	if err := app.Storage.Save(app.Ctx, app.Todos); err != nil {
		slog.Error("Failed to save todos on shutdown", slog.Any("error", err))
	}
}
