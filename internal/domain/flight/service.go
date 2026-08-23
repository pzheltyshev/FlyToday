package flight

import (
	"context"
	"time"
)

type FlightService struct {
	repo     Repository
	provider FlightProvider
}

func NewFlightService(repo Repository, provider FlightProvider) *FlightService {
	return &FlightService{
		repo:     repo,
		provider: provider,
	}
}

func (f *FlightService) SearchFlight(ctx context.Context, departureAirport Airport, arrivalAirport Airport, date time.Time) Flight {

	return Flight{}

}

func (f *FlightService) GetAirportByIAITCode(ctx context.Context, code string) Airport {

	airport, err := f.repo.GetAirportByIAITCode(ctx, code)

	if err != nil {
		return Airport{}
	}

	return airport

}
