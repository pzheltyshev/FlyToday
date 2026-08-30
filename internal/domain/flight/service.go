package flight

import (
	"context"
	"time"
	"log/slog"
)

type FlightService struct {
	repo     	Repository
	provider 	FlightProvider
	logger 		slog.Logger
}

func NewFlightService(repo Repository, provider FlightProvider, logger *slog.Logger) *FlightService {
	return &FlightService{
		repo:     repo,
		provider: provider,
		logger: logger,
	}
}

func (f *FlightService) SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) Flight {

	flightSegments, err := f.provider.SearchFlight(ctx, departureAirport, arrivalAirport, date)

	if err != nil {
		return Flight{}
	}

	f.repo.SaveRequest(flightSegments)

	flight := ConvSegmentsRawToFlight(&flightSegments)

	return flight

}

func (f *FlightService) GetAirportByIAITCode(ctx context.Context, code string) Airport {

	airport, err := f.repo.GetAirportByIAITCode(ctx, code)

	if err != nil {
		return Airport{}
	}

	return airport

}

func 