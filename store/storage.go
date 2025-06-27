package store

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"todo/logutil"
)

type Storage[T any] struct {
	FileName string
}

func NewStorage[T any](fileName string) *Storage[T] {
	return &Storage[T]{FileName: fileName}
}

func (s *Storage[T]) Save(ctx context.Context, data T) error {
	fileData, err := json.MarshalIndent(data, "", "    ")

	if err != nil {
		slog.Error("Failed to marshal data", slog.Any("error", err))
		return err
	}

	err = os.WriteFile(s.FileName, fileData, 0644)
	if err != nil {
		slog.Error("Failed to write file", slog.String("file", s.FileName), slog.Any("error", err))
		return err
	}

	logutil.Logger(ctx).Info("Data saved", "file", s.FileName)

	slog.Info("Data saved successfully", slog.String("file", s.FileName))
	return nil
}

func (s *Storage[T]) Load(ctx context.Context, data *T) error {
	fileData, err := os.ReadFile(s.FileName)

	// fileData, err := os.ReadFile(s.FileName)
	if err != nil {
		slog.Error("Failed to read file", slog.String("file", s.FileName), slog.Any("error", err))
		return err
	}

	err = json.Unmarshal(fileData, data)
	if err != nil {
		slog.Error("Failed to unmarshal data", slog.Any("error", err))
		return err
	}

	logutil.Logger(ctx).Info("Data loaded", "file", s.FileName)

	slog.Info("Data loaded successfully", slog.String("file", s.FileName))
	return nil

	// return json.Unmarshal(fileData, data)
}
