package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"html/template"
	"path/filepath"
	"todo/logutil"
	"todo/store"

	"github.com/google/uuid"
)

type App struct {
	Actor   *store.TodoActor
	Storage *store.Storage[store.Todos]
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func NewApp() *App {
	traceID := fmt.Sprintf("trace-%d", os.Getpid())
	ctx, cancel := SignalNotifyContext()
	ctx = logutil.WithTraceID(ctx, traceID)

	actor := store.NewTodoActor()

	storage := store.NewStorage[store.Todos]("todos.json")
	var loaded store.Todos
	if err := storage.Load(ctx, &loaded); err != nil {
		slog.Warn("Could not load todos", slog.Any("error", err))
	} else {
		for _, t := range loaded {
			actor.Add(ctx, t.Description)
		}
	}

	return &App{
		Actor:   actor,
		Storage: storage,
		Ctx:     ctx,
		Cancel:  cancel,
	}
}

func (a *App) Shutdown() {
	slog.Info("Shutting down and saving todos...")
	defer a.Cancel()

	if err := a.Storage.Save(a.Ctx, a.Actor.GetAll()); err != nil {
		slog.Error("Failed to save todos", slog.Any("error", err))
	}
}

func SignalNotifyContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		slog.Info("Interrupt received")
		cancel()
	}()

	return ctx, cancel
}

func (a *App) StartHTTPServer(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/about.html")
	})

	mux.HandleFunc("/list", withTrace(a.handleList))

	mux.HandleFunc("/create", withTrace(a.handleCreate))
	mux.HandleFunc("/get", withTrace(a.handleGet))
	mux.HandleFunc("/update", withTrace(a.handleUpdate))
	mux.HandleFunc("/delete", withTrace(a.handleDelete))

	slog.Info("Starting HTTP server on port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("HTTP server failed", slog.Any("error", err))
	}
}

func (a *App) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err := a.Actor.Add(ctx, body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.Storage.Save(ctx, a.Actor.GetAll())

	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (a *App) handleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	index, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo := a.Actor.Get(index)
	json.NewEncoder(w).Encode(todo)
}

func (a *App) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err := a.Actor.Edit(ctx, body.ID, body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.Storage.Save(ctx, a.Actor.GetAll())
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (a *App) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	index, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	todos := a.Actor.GetAll()

	if index < 0 || index >= len(todos) {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := a.Actor.Delete(ctx, index); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Storage.Save(ctx, a.Actor.GetAll()); err != nil {
		http.Error(w, "failed to persist changes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func withTrace(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.NewString()
		ctx := logutil.WithTraceID(r.Context(), traceID)

		logutil.Logger(ctx).Info("Received request", "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (a *App) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tmplPath := filepath.Join("web", "list.html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logutil.Logger(ctx).Error("template parse error", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, a.Actor.GetAll())
	if err != nil {
		logutil.Logger(ctx).Error("template exec error", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
