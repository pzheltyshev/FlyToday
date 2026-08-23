package app

import (
	"database/sql"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
	"github.com/pzheltyshev/FlyToday/internal/provider/mock"
	handler "github.com/pzheltyshev/FlyToday/internal/transport/http"
)

func BuildDependencies(db *sql.DB) {

	flightRepo := postgres.NewFlightRepository(db)

	flightProvider := mock.NewFlightProvider()

	flightService := flight.NewFlightService(flightRepo, flightProvider)

	flightHandler := handler.NewFlightHandler(flightService)

	flightHandler.Init()
}
