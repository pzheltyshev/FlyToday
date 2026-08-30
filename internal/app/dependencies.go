package app

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
	"github.com/pzheltyshev/FlyToday/internal/provider/mock"
	handler "github.com/pzheltyshev/FlyToday/internal/transport/http"
)

func newLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	return slog.New(handler)
}

func BuildDependencies(db *sql.DB) {

	logger := newLogger()

	flightRepo := postgres.NewFlightRepository(db)

	flightProvider := mock.NewFlightProvider()

	flightService := flight.NewFlightService(flightRepo, flightProvider, logger)

	flightHandler := handler.NewFlightHandler(flightService)

	flightHandler.Init()
}
