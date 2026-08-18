package flight

import (
	"context"
	"time"
)

type FlightService struct {
	repo Repository
}

func NewFlightService(repo Repository) *FlightService {
	return &FlightService{
		repo: repo,
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
