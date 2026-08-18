package flight

import (
	"context"
	"time"
)

type Repository interface {
	SearchFlight(ctx context.Context, departureAirport Airport, arrival Airport, date time.Time) (Flight, error)
	GetAirportByIAITCode(ctx context.Context, code string) (Airport, error)
}
