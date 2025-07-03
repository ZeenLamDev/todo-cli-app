package main

import (
	"log/slog"
	"os"

	"todo/app"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	application := app.NewApp()

	go application.StartHTTPServer("8080")
	slog.Info("🚀 App is running on http://localhost:8080 ... Ctrl+C to stop")
	<-application.Ctx.Done()

	slog.Info("🛑 Interrupt received. Shutting down...")
	application.Shutdown()
}
